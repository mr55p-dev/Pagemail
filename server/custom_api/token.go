package custom_api

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenData struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Data TokenData `json:"data"`
}

func getKeyFromEnvironment(env_var string) ([]byte, error) {
	key_file, is_set := os.LookupEnv(env_var)
	if len(key_file) == 0 ||!is_set {
		return nil, fmt.Errorf("Failed to obtain token signing key")
	}

	key_file_path, err := filepath.Abs(key_file)
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain token signing key: %s", err)
	}

	key, err := os.ReadFile(key_file_path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read token signing key data: %s", err)
	}

	return key, nil
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	key, err := getKeyFromEnvironment("PAGEMAIL_TOKEN_SIGNING_PRIVATE_KEY")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(key)
}

func getPublicKey() (*rsa.PublicKey, error) {
	key, err := getKeyFromEnvironment("PAGEMAIL_TOKEN_SIGNING_PUBLIC_KEY")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(key)
}

func NewTokenRoute(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		user_record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if user_record == nil {
			return apis.NewForbiddenError("User not found or not authenticated", nil)
		}

		token_id := uuid.NewString()
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, TokenClaims{
			user_record.GetId(),
			jwt.RegisteredClaims{
				Issuer:   "pagemail.io",
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Subject:  "pageApiAdd",
				ID:       token_id,
			},
		})

		key, err := getPrivateKey()
		if err != nil {
			return apis.NewApiError(500, fmt.Sprintf("Could not fetch token signing key: %s", err), nil)
		}
		tkn, err := token.SignedString(key)
		if err != nil {
			return apis.NewApiError(500, fmt.Sprintf("API failed to generate new token: %s", err), nil)
		}

		form := forms.NewRecordUpsert(app, user_record)
		form.LoadData(map[string]any{
			"userTokenId": token_id,
		})
		if err := form.Submit(); err != nil {
			return apis.NewApiError(500, fmt.Sprintf("API failed to generate new token: %s", err), nil)
		}

		return c.JSON(200, TokenResponse{Data: TokenData{Token: tkn}})
	}
}

func VerifyTokenMiddleware(app *pocketbase.PocketBase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token_string := c.Request().Header.Get("Authorization")
			token, err := jwt.ParseWithClaims(token_string, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) { return getPublicKey() })

			if err != nil {
				return apis.NewForbiddenError("Provided token is not valid", nil)
			}

			claims, ok := token.Claims.(*TokenClaims)
			if !ok || !token.Valid {
				return apis.NewForbiddenError("Provided token is not valid", nil)
			}

			user_id := claims.UserID
			user_record, err := app.Dao().FindRecordById("users", user_id)
			if err != nil {
				return apis.NewForbiddenError(fmt.Sprintf("Failed to fetch user: %s", err), nil)
			}
			if user_record.Get("userTokenId") != claims.ID {
				return apis.NewForbiddenError("Token is expired", nil)
			}

			c.Set("TokenClaims", claims)

			return next(c)
		}
	}
}
