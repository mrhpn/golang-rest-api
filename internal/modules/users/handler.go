package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

// Handler handles user-related HTTP endpoints such as user profile access and account management operations.
type Handler struct {
	userService Service
}

// NewHandler constructs a users Handler with its required service dependency.
func NewHandler(userService Service) *Handler {
	return &Handler{userService: userService}
}

// Create user godoc
//
//	@Summary		Create user
//	@Description	Create a user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createUserRequest	true	"CreateUserRequest"
//	@Success		201		{object}	users.UserResponse
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Failure		401		{object}	httpx.ErrorResponse
//	@Failure		500		{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users [post]
func (h *Handler) Create(c *gin.Context) {
	var req createUserRequest
	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.userService.Create(httpx.ReqCtx(c), req)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusCreated,
		UserResponse{ID: user.ID, Email: user.Email, Role: user.Role},
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
	var params iDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	user, err := h.userService.GetByID(httpx.ReqCtx(c), params.ID)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(
		c,
		http.StatusOK,
		UserResponse{ID: user.ID, Email: user.Email, Role: user.Role},
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
	var params iDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.userService.Delete(httpx.ReqCtx(c), params.ID); err != nil {
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
	var params iDParam

	if err := httpx.BindURI(c, &params); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	if err := h.userService.Restore(httpx.ReqCtx(c), params.ID); err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusOK, nil)
}
