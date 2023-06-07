package custom_api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/search"
	"golang.org/x/oauth2"
)

func HandleGoogleLogin(token oauth2.Token, form *forms.RecordOAuth2Login) models.Record {
	provider := auth.NewGoogleProvider()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	provider.SetContext(ctx)

	// fetch external auth user
	authUser, err := provider.FetchAuthUser(&token)
	if err != nil {
		return err
	}

	var authRecord *models.Record

	// check for existing relation with the auth record
	rel, _ := form.dao.FindExternalAuthByProvider(form.Provider, authUser.Id)
	switch {
	case rel != nil:
		authRecord, err = form.dao.FindRecordById(form.collection.Id, rel.RecordId)
		if err != nil {
			return nil, authUser, err
		}
	case form.loggedAuthRecord != nil && form.loggedAuthRecord.Collection().Id == form.collection.Id:
		// fallback to the logged auth record (if any)
		authRecord = form.loggedAuthRecord
	case authUser.Email != "":
		// look for an existing auth record by the external auth record's email
		authRecord, _ = form.dao.FindAuthRecordByEmail(form.collection.Id, authUser.Email)
	}

	return *authRecord
}

func GAuthCustomBindApp(app *pocketbase.PocketBase) echo.HandlerFunc {

	return func(c echo.Context) error {
		collection, _ := c.Get("collection").(*models.Collection)
		if collection == nil {
			return c.String(http.StatusNotFound, "Missing collection context.")
		}

		if !collection.AuthOptions().AllowOAuth2Auth {
			return c.String(http.StatusBadRequest, "The collection is not configured to allow OAuth2 authentication.")
		}

		var fallbackAuthRecord *models.Record

		loggedAuthRecord, _ := c.Get("authRecord").(*models.Record)
		if loggedAuthRecord != nil && loggedAuthRecord.Collection().Id == collection.Id {
			fallbackAuthRecord = loggedAuthRecord
		}

		form := forms.NewRecordOAuth2Login(app, collection, fallbackAuthRecord)
		if readErr := c.Bind(form); readErr != nil {
			return c.String(http.StatusInternalServerError, "An error occurred while loading the submitted data.")
		}

		event := new(core.RecordAuthWithOAuth2Event)
		event.HttpContext = c
		event.Collection = collection
		event.ProviderName = form.Provider
		event.IsNewRecord = false

		form.SetBeforeNewRecordCreateFunc(func(createForm *forms.RecordUpsert, authRecord *models.Record, authUser *auth.AuthUser) error {
			return createForm.DrySubmit(func(txDao *daos.Dao) error {
				event.IsNewRecord = true
				// clone the current request data and assign the form create data as its body data
				requestData := *apis.RequestData(c)
				requestData.Data = form.CreateData

				createRuleFunc := func(q *dbx.SelectQuery) error {
					admin, _ := c.Get("admin").(*models.Admin)
					if admin != nil {
						return nil // either admin or the rule is empty
					}

					if collection.CreateRule == nil {
						return errors.New("Only admins can create new accounts with OAuth2")
					}

					if *collection.CreateRule != "" {
						resolver := resolvers.NewRecordFieldResolver(txDao, collection, &requestData, true)
						expr, err := search.FilterData(*collection.CreateRule).BuildExpr(resolver)
						if err != nil {
							return err
						}
						resolver.UpdateQuery(q)
						q.AndWhere(expr)
					}

					return nil
				}

				if _, err := txDao.FindRecordById(collection.Id, createForm.Id, createRuleFunc); err != nil {
					return fmt.Errorf("Failed create rule constraint: %w", err)
				}

				return nil
			})
		})

		_, _, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*forms.RecordOAuth2LoginData]) forms.InterceptorNextFunc[*forms.RecordOAuth2LoginData] {
			return func(data *forms.RecordOAuth2LoginData) error {
				event.Record = data.Record
				event.OAuth2User = data.OAuth2User
				event.ProviderClient = data.ProviderClient

				return app.OnRecordBeforeAuthWithOAuth2Request().Trigger(event, func(e *core.RecordAuthWithOAuth2Event) error {
					data.Record = e.Record
					data.OAuth2User = e.OAuth2User

					if err := next(data); err != nil {
						return NewBadRequestError("Failed to authenticate.", err)
					}

					e.Record = data.Record
					e.OAuth2User = data.OAuth2User

					meta := struct {
						*auth.AuthUser
						IsNew bool `json:"isNew"`
					}{
						AuthUser: e.OAuth2User,
						IsNew:    event.IsNewRecord,
					}

					return RecordAuthResponse(app, e.HttpContext, e.Record, meta)
				})
			}
		})

		if submitErr == nil {
			if err := app.OnRecordAfterAuthWithOAuth2Request().Trigger(event); err != nil && api.app.IsDebug() {
				log.Println(err)
			}
		}

		return submitErr
	}
}
