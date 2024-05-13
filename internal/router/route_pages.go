package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/preview"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

type GetPagesRequest struct {
	Page string `query:"p"`
}

type GetPageRequest struct {
	PageID string `param:"page_id"`
}

type PostPageRequest struct {
	Url string `form:"url"`
}

func (router *Router) GetPages(w http.ResponseWriter, r *http.Request) {
	req := requestBind[GetPagesRequest](w, r)
	if req == nil {
		return
	}

	user := dbqueries.GetUser(r.Context())
	page, err := strconv.Atoi(req.Page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// pageinate for fun
	pages = pages[(page-1)*render.PAGE_SIZE : page*render.PAGE_SIZE]
	staticRender(render.PageList(pages, page), w, r)
}

func (router *Router) DeletePages(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())

	n, err := router.DBClient.DeletePagesByUserId(r.Context(), user.ID)
	if err != nil {
		genericResponse(w, http.StatusInternalServerError)
		return
	}
	staticRender(render.SavePageSuccess(fmt.Sprintf("Deleted %d pages", n)), w, r)

}

func (router *Router) GetPage(w http.ResponseWriter, r *http.Request) {
	req := requestBind[GetPageRequest](w, r)
	if req == nil {
		return
	}

	user := dbqueries.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), req.PageID)
	if err != nil {
		http.Error(w, "Failed to get page id", http.StatusInternalServerError)
		return
	}
	if page.UserID != user.ID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	_, err = router.DBClient.DeletePageById(r.Context(), req.PageID)
	if err != nil {
		genericResponse(w, http.StatusInternalServerError)
		return
	}
	staticRender(render.PageCard(&page), w, r)
}

func (router *Router) PostPage(w http.ResponseWriter, r *http.Request) {
	req := requestBind[PostPageRequest](w, r)
	if req == nil {
		return
	}

	user := dbqueries.GetUser(r.Context())
	url := req.Url
	if url == "" {
		staticRender(render.SavePageError("URL field must be present"), w, r)
		return
	}

	now := time.Now()
	page := dbqueries.Page{
		ID:      tools.GenerateNewId(20),
		UserID:  user.ID,
		Url:     url,
		Created: now,
		Updated: now,
	}
	err := router.DBClient.CreatePage(r.Context(), dbqueries.CreatePageParams{
		ID:      page.ID,
		UserID:  page.UserID,
		Url:     page.Url,
		Created: page.Created,
		Updated: page.Updated,
	})
	if err != nil {
		staticRender(render.SavePageError(err.Error()), w, r)
	}

	go func(cli *dbqueries.Queries) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := preview.FetchPreview(ctx, &page)
		if err != nil {
			return
		}
		err = cli.UpsertPage(ctx, dbqueries.UpsertPageParams{
			ID:                  page.ID,
			UserID:              page.UserID,
			Url:                 page.Url,
			Title:               page.Title,
			Description:         page.Description,
			ImageUrl:            page.ImageUrl,
			ReadabilityStatus:   page.ReadabilityStatus,
			ReadabilityTaskData: page.ReadabilityTaskData,
			IsReadable:          page.IsReadable,
			Created:             page.Created,
			Updated:             page.Updated,
		})
		if err != nil {
			return
		}
	}(router.DBClient)

	staticRender(render.PageCard(&page), w, r)
	return
}

func (router *Router) DeletePage(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), r.PathValue("page_id"))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if !router.Authorizer.ValUserAgainstPage(user.ID, page.UserID) {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	router.DBClient.DeletePageById(r.Context(), page.ID)
	w.WriteHeader(http.StatusOK)
}
