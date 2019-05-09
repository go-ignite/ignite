package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/logger"
	"github.com/go-ignite/ignite/models"
)

type UserHandler struct {
	*Handler
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		Handler: New(logger.GetUserHandlerLogger()),
	}
}

func (uh *UserHandler) Login(c *gin.Context) {
	req := new(api.UserLoginRequest)
	if err := c.ShouldBind(req); err != nil {
		uh.ErrJSON(c, http.StatusBadRequest, err)
		return
	}
	uh.logger.WithField("req", fmt.Sprintf("%+v", req)).Debug("Login")

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "generate hashed password error"))
		return
	}
	user, err := models.GetUserByNameAndPassword(req.Username, hashedPass)
	if err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get user error"))
		return
	}
	if user == nil {
		uh.ErrJSON(c, http.StatusUnauthorized, nil)
		return
	}

	token, err := uh.createToken(int64(user.ID))
	if err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, api.NewUserLoginResponse(token))
}

// 1 -> invite code not found
func (uh *UserHandler) Register(c *gin.Context) {
	req := new(api.UserRegisterRequest)
	if err := c.ShouldBind(req); err != nil {
		uh.ErrJSON(c, http.StatusBadRequest, err)
		return
	}
	uh.logger.WithField("req", fmt.Sprintf("%+v", req)).Debug("Register")

	ic, err := models.GetInviteCodeByCode(req.InviteCode)
	if err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get invite code error"))
		return
	}
	if ic == nil {
		uh.ErrJSON(c, http.StatusBadRequest, errors.Wrap(err, "invite code not found"), 1)
		return
	}

	uh.logger.WithField("inviteCodeID", ic.ID).Debug("Register")

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "generate hashed password error"))
		return
	}
	user := models.NewUser(req.Username, hashedPass, ic.ID)
	if err := user.Create(ic); err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "create user error"))
		return
	}

	token, err := uh.createToken(int64(user.ID))
	if err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, api.NewUserRegisterResponse(token))
}

func (uh *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetInt("id")

	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get user error"))
		return
	}
	if user == nil {
		uh.ErrJSON(c, http.StatusNotFound, nil)
		return
	}
	resp := new(api.User)
	if err := uh.copy(resp, user); err != nil {
		uh.ErrJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
