package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/pkg/request"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

type GetPagesRequest struct {
	Page string `query:"p"`
}

func (router *Router) GetPages(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := request.BindRequest[GetPagesRequest](w, r)
	if req == nil {
		return
	}

	user := auth.GetUser(r.Context())
	page, err := strconv.Atoi(req.Page)
	if err != nil {
		response.Error(w, r, pmerror.ErrBadPagination)
		return
	}

	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.ID)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to load users pages")
		response.Error(w, r, pmerror.NewInternalError("Failed to load your pages"))
		return
	}
	// pageinate for fun
	pages = pages[(page-1)*render.PAGE_SIZE : page*render.PAGE_SIZE]
	response.Component(render.PageList(pages, page), w, r)
}

func (router *Router) DeletePages(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	user := auth.GetUser(r.Context())

	n, err := router.DBClient.DeletePagesByUserId(r.Context(), user.ID)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to delete all pages")
		response.Error(w, r, pmerror.NewInternalError("Failed to delete pages"))
		return
	}
	response.Success(fmt.Sprintf("Deleted %d pages", n), w, r)
}

func (router *Router) GetPage(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	id := r.PathValue("page_id")
	if id == "" {
		response.Error(w, r, pmerror.ErrNoParam)
		return
	}
	user := auth.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), id)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to load page")
		response.Error(w, r, pmerror.NewInternalError("Failed to get page"))
		return
	}
	if page.UserID != user.ID {
		response.Error(w, r, pmerror.ErrNoPage)
		return
	}
	if request.IsHtmx(r) {
		if page.PreviewState != "unknown" {
			w.WriteHeader(286) // tell htmx to stop polling
		}
	}
	response.Component(render.PageCard(&page), w, r)
}

type PostPageRequest struct {
	Url string `form:"url"`
}

func (router *Router) PostPage(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	req := request.BindRequest[PostPageRequest](w, r)
	if req == nil {
		return
	}

	user := auth.GetUser(r.Context())
	url := req.Url
	if url == "" {
		response.Error(w, r, pmerror.ErrNoParam)
		return
	}

	now := time.Now()
	page := dbqueries.Page{
		ID:      tools.GenerateNewId(20),
		UserID:  user.ID,
		Url:     url,
		PreviewState: "unknown",
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
		response.Error(w, r, pmerror.NewInternalError("Failed to save page"))
		return
	}

	router.Previewer.Queue(page.ID)

	if request.IsHtmx(r) {
		response.Component(render.PageCard(&page), w, r)
	} else {
		response.Success("Added page successfully", w, r)
	}
	return
}

func (router *Router) DeletePage(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	user := auth.GetUser(r.Context())
	page, err := router.DBClient.ReadPageById(r.Context(), r.PathValue("page_id"))
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to find page in db")
		response.Error(w, r, pmerror.NewInternalError("Failed to delete page"))
		return
	}

	if !(user.ID == page.UserID) {
		logger.WithError(err).ErrorCtx(r.Context(), "Attempt to delete another users page")
		response.Error(w, r, pmerror.ErrNotAllowed)
		return
	}

	_, err = router.DBClient.DeletePageById(r.Context(), page.ID)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to delete page from db")
		response.Error(w, r, pmerror.NewInternalError("Failed to delete page"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
