package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	userpb "github.com/s-usmonalizoda25/protoCinemaService/gen/user"
)

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
