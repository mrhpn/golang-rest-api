package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/pagination"
)

type Service interface {
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, opts *pagination.QueryOptions) ([]*User, *httpx.PaginationMeta, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	Block(ctx context.Context, id string) error
	Reactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) (*User, error)
}

// Handler handles user-related HTTP endpoints such as user profile access and account management operations.
type Handler struct {
	service Service
}

// NewHandler constructs a users Handler with its required service dependency.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create user godoc
//
//	@Summary		Create user
//	@Description	Create a user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		users.CreateUserRequest	true	"CreateUserRequest"
//	@Success		201		{object}	users.UserResponse
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Failure		401		{object}	httpx.ErrorResponse
//	@Failure		500		{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := httpx.BindAndValidateJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.Create(httpx.ReqCtx(c), req)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusCreated,
		ToUserResponse(user),
	)
}

// Get a user godoc
//
//	@Summary		Get user
//	@Description	Get a user by their ULID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		201	{object}	users.UserResponse
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	var params IDParam
	if err := httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.GetByID(httpx.ReqCtx(c), params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusOK,
		ToUserResponse(user),
	)
}

// List users godoc
//
//	@Summary		List users
//	@Description	Get a paginated list of users with search and sorting
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int		false	"Page number (default: 1)"					default(1)	minimum(1)
//	@Param			limit			query		int		false	"Items per page (default: 10, max: 100)"		default(10)	minimum(1)
//	@Param			search			query		string	false	"Search text (case-insensitive)"
//	@Param			search_columns	query		[]string	false	"Columns to search in (default: all searchable)"
//	@Param			sort_by			query		string	false	"Field to sort by (email, role, created_at)"
//	@Param			order			query		string	false	"Sort order (asc or desc)"					Enums(asc, desc)	default(desc)
//	@Param			exact_match		query		bool	false	"Use exact match for search (default: false)"
//	@Success		200				{object}	httpx.SuccessResponse{data=[]users.UserResponse,meta=httpx.PaginationMeta}
//	@Failure		400				{object}	httpx.ErrorResponse
//	@Failure		401				{object}	httpx.ErrorResponse
//	@Failure		500				{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users [get]
func (h *Handler) List(c *gin.Context) {
	var query pagination.QueryList
	if err := httpx.BindAndValidateQuery(c, &query); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	opts := pagination.NewQueryOptions(
		&query,
		pagination.SortSearchPolicy{
			SortableCols:   []string{"role", "created_at", "updated_at"},
			SearchableCols: []string{"email"},
		},
	)

	users, meta, err := h.service.List(httpx.ReqCtx(c), opts)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OKWithMeta(
		c,
		http.StatusOK,
		ToUserResponseList(users),
		meta,
	)
}

// Delete a user godoc
//
//	@Summary		Delete user
//	@Description	Delete a user by their ULID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	var params IDParam
	if err := httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Delete(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Restore a deleted user godoc
//
//	@Summary		Restore user
//	@Description	Restore a deleted user by their ULID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/{id}/restore [put]
func (h *Handler) Restore(c *gin.Context) {
	var params IDParam
	if err := httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Restore(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, nil)
}

// Block a user godoc
//
//	@Summary		Block user
//	@Description	Block a user by their ULID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200	{object}	users.UserResponse
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/{id}/block [put]
func (h *Handler) Block(c *gin.Context) {
	var params IDParam
	if err := httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Block(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.GetByID(httpx.ReqCtx(c), params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, ToUserResponse(user))
}

// Reactivate a user godoc
//
//	@Summary		Reactivate user
//	@Description	Reactivate a user by their ULID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200	{object}	users.UserResponse
//	@Failure		400	{object}	httpx.ErrorResponse
//	@Failure		401	{object}	httpx.ErrorResponse
//	@Failure		500	{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/{id}/reactivate [put]
func (h *Handler) Reactivate(c *gin.Context) {
	var params IDParam
	if err := httpx.BindAndValidateURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.service.Reactivate(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.service.GetByID(httpx.ReqCtx(c), params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, ToUserResponse(user))
}
