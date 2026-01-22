package posts

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/middlewares"
	"github.com/mrhpn/go-rest-api/internal/pagination"
)

// PostService defines the business logic for managing posts.
type PostService interface {
	Create(ctx context.Context, userID string, req CreatePostRequest) (*Post, error)
	GetByID(ctx context.Context, id string) (*Post, error)
	GetByUserID(ctx context.Context, userID string, opts *pagination.QueryOptions) ([]*Post, *httpx.PaginationMeta, error)
	List(ctx context.Context, opts *pagination.QueryOptions) ([]*Post, *httpx.PaginationMeta, error)
	Update(ctx context.Context, id string, userID string, req UpdatePostRequest) error
	Delete(ctx context.Context, id string, userID string) error
}

// Handler handles post-related HTTP endpoints such as post creation, reading, updating, and deletion.
type Handler struct {
	service PostService
}

// NewHandler constructs a posts Handler with its required service dependency.
func NewHandler(service PostService) *Handler {
	return &Handler{service: service}
}

// Create post godoc
//
//	@Summary		Create post
//	@Description	Create a new post
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			request	body		posts.CreatePostRequest	true	"CreatePostRequest"
//	@Success		201		{object}	posts.PostResponse
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Failure		401		{object}	httpx.ErrorResponse
//	@Failure		500		{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/posts [post]
func (h *Handler) Create(c *gin.Context) {
	// Get user from context
	user, err := middlewares.GetUser(httpx.ReqCtx(c))
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	var req CreatePostRequest
	if err = httpx.BindAndValidateJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	post, err := h.service.Create(httpx.ReqCtx(c), user.UserID, req)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusCreated,
		ToPostResponse(post),
	)
}

// Get post godoc
//
//	@Summary		Get post
//	@Description	Get a post by its ID
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	posts.PostResponse
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		404	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/posts/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	var params IDParam
	if err := httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	post, err := h.service.GetByID(httpx.ReqCtx(c), params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusOK,
		ToPostResponse(post),
	)
}

// List posts godoc
//
//	@Summary		List posts
//	@Description	Get a paginated list of posts with search and sorting
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default: 1)"					default(1)	minimum(1)
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"		default(10)	minimum(1)
//	@Param			search		query		string	false	"Search text (case-insensitive)"
//	@Param			sort_by		query		string	false	"Field to sort by (title, created_at, updated_at)"
//	@Param			order		query		string	false	"Sort order (asc or desc)"					Enums(asc, desc)	default(desc)
//	@Success		200			{object}	httpx.SuccessResponse{data=[]posts.PostResponse,meta=httpx.PaginationMeta}
//	@Failure		400			{object}	httpx.ErrorResponse
//	@Failure		401			{object}	httpx.ErrorResponse
//	@Failure		500			{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/posts [get]
func (h *Handler) List(c *gin.Context) {
	var query pagination.QueryList
	if err := httpx.BindAndValidateQuery(c, &query); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	opts := pagination.NewQueryOptions(
		&query,
		pagination.SortSearchPolicy{
			SortableCols:   []string{"title", "created_at", "updated_at"},
			SearchableCols: []string{"title", "content"},
		},
	)

	posts, meta, err := h.service.List(httpx.ReqCtx(c), opts)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OKWithMeta(
		c,
		http.StatusOK,
		ToPostResponseList(posts),
		meta,
	)
}

// ListMyPosts godoc
//
//	@Summary		List my posts
//	@Description	Get a paginated list of posts created by the authenticated user
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default: 1)"					default(1)	minimum(1)
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"		default(10)	minimum(1)
//	@Param			search		query		string	false	"Search text (case-insensitive)"
//	@Param			sort_by		query		string	false	"Field to sort by (title, created_at, updated_at)"
//	@Param			order		query		string	false	"Sort order (asc or desc)"					Enums(asc, desc)	default(desc)
//	@Success		200			{object}	httpx.SuccessResponse{data=[]posts.PostResponse,meta=httpx.PaginationMeta}
//	@Failure		400			{object}	httpx.ErrorResponse
//	@Failure		401			{object}	httpx.ErrorResponse
//	@Failure		500			{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/posts/my [get]
func (h *Handler) ListMyPosts(c *gin.Context) {
	// Get user from context
	user, err := middlewares.GetUser(httpx.ReqCtx(c))
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	var query pagination.QueryList
	if err = httpx.BindAndValidateQuery(c, &query); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	opts := pagination.NewQueryOptions(
		&query,
		pagination.SortSearchPolicy{
			SortableCols:   []string{"title", "created_at", "updated_at"},
			SearchableCols: []string{"title", "content"},
		},
	)

	posts, meta, err := h.service.GetByUserID(httpx.ReqCtx(c), user.UserID, opts)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OKWithMeta(
		c,
		http.StatusOK,
		ToPostResponseList(posts),
		meta,
	)
}

// Update post godoc
//
//	@Summary		Update post
//	@Description	Update a post by its ID (only the owner can update)
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Post ID"
//	@Param			request	body		posts.UpdatePostRequest	true	"UpdatePostRequest"
//	@Success		200		{object}	posts.PostResponse
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Failure		401		{object}	httpx.ErrorResponse
//	@Failure		403		{object}	httpx.ErrorResponse
//	@Failure		404		{object}	httpx.ErrorResponse
//	@Failure		500		{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/posts/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	// Get user from context
	user, err := middlewares.GetUser(httpx.ReqCtx(c))
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	var params IDParam
	if err = httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	var req UpdatePostRequest
	if err = httpx.BindAndValidateJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err = h.service.Update(httpx.ReqCtx(c), params.ID, user.UserID, req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	// Fetch updated post
	post, err := h.service.GetByID(httpx.ReqCtx(c), params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, ToPostResponse(post))
}

// Delete post godoc
//
//	@Summary		Delete post
//	@Description	Delete a post by its ID (only the owner can delete)
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Post ID"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		403	{object}	httpx.ErrorResponse
//	@Failure		404	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/posts/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	// Get user from context
	user, err := middlewares.GetUser(httpx.ReqCtx(c))
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	var params IDParam
	if err = httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err = h.service.Delete(httpx.ReqCtx(c), params.ID, user.UserID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
