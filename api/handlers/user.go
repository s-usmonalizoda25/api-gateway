package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/models"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	userpb "github.com/s-usmonalizoda25/protoCinemaService/gen/user"
)

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

	c.JSON(http.StatusOK, response)
}
