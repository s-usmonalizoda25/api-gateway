package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/models"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	userpb "github.com/s-usmonalizoda25/protoCinemaService/gen/user"
)

// Register
// @Summary Register a new user
// @Description Register a new user with name, email, password, etc.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User Registration Info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/register [post]
func (h *handler) Register(c *gin.Context) {
	var body models.CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		errs.HandleValidationError(c, err)
		return
	}

	response, err := h.serviceManager.UserService().Add(c.Request.Context(), &userpb.CreateUserRequest{
		Name:     body.Username,
		Email:    body.Email,
		Password: body.Password,
		Age:      body.Age,
		Phone:    body.Phone,
	})

	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedRegister, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser
// @Summary Get User
// @Description Get user by id
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/{user_id} [get]
func (h *handler) GetUser(c *gin.Context) {
	idStr := c.Param("user_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	response, err := h.serviceManager.UserService().GetByID(c.Request.Context(), &userpb.GetUserRequest{
		Id: id,
	})

	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetUser, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Login
// @Summary User Login
// @Description Authenticate user and get tokens
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/user/login [post]
func (h *handler) Login(c *gin.Context) {
	var body models.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		errs.HandleValidationError(c, err)
		return
	}

	response, err := h.serviceManager.UserService().Login(c.Request.Context(), &userpb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedLogin, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            response.Id,
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
	})
}

// GetMyProfile
// @Summary Get My Profile
// @Description Get profile details of the currently authenticated user
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/user/me [get]
func (h *handler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errs.HandleAuthError(c, h.log, errs.MsgUnauthorized)
		return
	}

	uid := int64(userID.(float64))

	response, err := h.serviceManager.UserService().GetByID(c.Request.Context(), &userpb.GetUserRequest{
		Id: uid,
	})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetUser, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateMyProfile
// @Summary Update My Profile
// @Description Update profile details of the currently authenticated user
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateUserRequest true "Updated Profile Info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/user/me [put]	
func (h *handler) UpdateMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errs.HandleAuthError(c, h.log, errs.MsgUnauthorized)
		return
	}
	uid := int64(userID.(float64))

	var body models.UpdateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		errs.HandleValidationError(c, err)
		return
	}

	_, err := h.serviceManager.UserService().Update(c.Request.Context(), &userpb.UpdateUserRequest{
		Id:    uid,
		Name:  body.Username,
		Phone: body.Phone,
	})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedUpdateUser, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
