package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/logger"
)

type AdminHandler struct {
	*Handler
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		Handler: New(logger.GetAdminHandlerLogger()),
	}
}

func (ah *AdminHandler) createToken() (string, error) {
	return ah.Handler.createToken(-1)
}

func (ah *AdminHandler) Login(c *gin.Context) {
	req := new(api.AdminLoginRequest)
	if err := c.ShouldBind(req); err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}

	if !config.C.Admin.Match(req.Username, req.Password) {
		ah.ErrJSON(c, http.StatusUnauthorized, nil)
		return
	}

	token, err := ah.createToken()
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, api.NewAdminLoginResponse(token))
}
