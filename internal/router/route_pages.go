package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/preview"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

type GetPagesRequest struct {
	Page string `query:"p"`
}

func (router *Router) GetPages(w http.ResponseWriter, r *http.Request) {
	req := requestBind[GetPagesRequest](w, r)
	if req == nil {
		return
	}

	user := auth.GetUser(r.Context())
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
	componentRender(render.PageList(pages, page), w, r)
}

func (router *Router) DeletePages(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())

	n, err := router.DBClient.DeletePagesByUserId(r.Context(), user.ID)
	if err != nil {
		genericResponse(w, http.StatusInternalServerError)
		return
	}
	componentRender(render.SavePageSuccess(fmt.Sprintf("Deleted %d pages", n)), w, r)

}

func (router *Router) GetPage(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	id := r.PathValue("page_id")
	if id == "" {
		http.Error(w, "Missing page id", http.StatusBadRequest)
		return
	}
	user := auth.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), id)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to load page")
		http.Error(w, "Failed to get page id", http.StatusInternalServerError)
		return
	}
	if page.UserID != user.ID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if isHtmx(r) {
		if page.PreviewState != "unknown" {
			w.WriteHeader(286) // tell htmx to stop polling
		}
	}
	componentRender(render.PageCard(&page), w, r)
}

type PostPageRequest struct {
	Url string `form:"url"`
}

func (router *Router) PostPage(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := requestBind[PostPageRequest](w, r)
	if req == nil {
		return
	}

	user := auth.GetUser(r.Context())
	url := req.Url
	if url == "" {
		errorResponse(w, r, "Missing URL in request", http.StatusBadRequest)
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
		logger.WithError(err).ErrorCtx(r.Context(), "Error ")
		errorResponse(w, r, "Failed to save page", http.StatusInternalServerError)
		return
	}

	go func(cli *dbqueries.Queries, page dbqueries.Page) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := preview.FetchPreview(ctx, &page)
		pageUpdate := dbqueries.UpdatePagePreviewParams{
			ID:          page.ID,
			Title:       page.Title,
			Description: page.Description,
			ImageUrl:    page.ImageUrl,
			Updated:     time.Now(),
		}
		if err == nil {
			pageUpdate.PreviewState = "success"
		} else {
			pageUpdate.PreviewState = "error"
		}

		err = cli.UpdatePagePreview(ctx, pageUpdate)
		if err != nil {
			return
		}
	}(router.DBClient, page)

	if isHtmx(r) {
		componentRender(render.PageCard(&page), w, r)
	} else {
		textRender("Added page successfully", w, r)
	}
	return
}

func (router *Router) DeletePage(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	user := auth.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), r.PathValue("page_id"))
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to find page in db")
		errorResponse(w, r, "Failed to delete page", http.StatusInternalServerError)
		return
	}

	if !(user.ID == page.UserID) {
		logger.WithError(err).ErrorCtx(r.Context(), "Attempt to delete another users page")
		errorResponse(w, r, "Permission denied", http.StatusForbidden)
		return
	}

	_, err = router.DBClient.DeletePageById(r.Context(), page.ID)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to delete page from db")
		errorResponse(w, r, "Failed to delete page", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
