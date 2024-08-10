package router

import (
	"fmt"
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/render/components"
	"github.com/mr55p-dev/pagemail/internal/render/views"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/pkg/request"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

func (router *Router) GetDashboard(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	pages, err := queries.New(router.db).ReadPagesByUserId(r.Context(), queries.ReadPagesByUserIdParams{
		UserID: user.ID,
		Limit:  render.PAGE_SIZE,
		Offset: 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Component(views.Dashboard(user, pages), w, r)
}

type GetPagesRequest struct {
	Page string `query:"p"`
}

func (router *Router) DeletePages(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	user := auth.GetUser(r.Context())

	n, err := queries.New(router.db).DeletePagesByUserId(r.Context(), user.ID)
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
	page, err := queries.New(router.db).ReadPageById(r.Context(), id)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to load page")
		response.Error(w, r, pmerror.NewInternalError("Failed to get page"))
		return
	}
	if page.UserID != user.ID {
		response.Error(w, r, pmerror.ErrNoPage)
		return
	}
	response.Component(components.PageCard(&page), w, r)
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

	page := queries.Page{}
	page, err := queries.New(router.db).CreatePage(r.Context(), queries.CreatePageParams{
		ID:     tools.NewPageId(),
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Error ")
		response.Error(w, r, pmerror.NewInternalError("Failed to save page"))
		return
	}

	router.Previewer.Queue(page.ID)

	if request.IsHtmx(r) {
		response.Component(components.PageCard(&page), w, r)
	} else {
		response.Success("Added page successfully", w, r)
	}
	return
}

func (router *Router) DeletePage(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	user := auth.GetUser(r.Context())
	page, err := queries.New(router.db).ReadPageById(r.Context(), r.PathValue("page_id"))
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

	_, err = queries.New(router.db).DeletePageById(r.Context(), page.ID)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "Failed to delete page from db")
		response.Error(w, r, pmerror.NewInternalError("Failed to delete page"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
