package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/model"
)

func (s *Service) UserLogin(c *gin.Context) {
	req := new(api.UserLoginRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}
	user, err := s.opts.ModelHandler.GetUserByNameAndPassword(req.Username, hashedPass)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		s.errJSON(c, http.StatusUnauthorized, nil)
		return
	}

	token, err := s.createToken(user.ID)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &api.UserLoginResponse{Token: token})
}

// 1 -> invite code not found
func (s *Service) UserRegister(c *gin.Context) {
	req := new(api.UserRegisterRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	ic, err := s.opts.ModelHandler.GetInviteCode(req.InviteCode)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}
	if ic == nil {
		s.errJSON(c, http.StatusBadRequest, err, 1)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	user := model.NewUser(req.Username, hashedPass, ic.ID)
	if err := s.opts.ModelHandler.CreateUser(user); err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	token, err := s.createToken(user.ID)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, api.UserResisterResponse{Token: token})
}

func (s *Service) GetUserInfo(c *gin.Context) {
	userID := c.GetString("id")
	user, err := s.opts.ModelHandler.GetUserByID(userID)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		s.errJSON(c, http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, &api.User{
		ID:   userID,
		Name: user.Name,
	})
}

func (s *Service) Sync(c *gin.Context) {
	//for {
	//	nam := s.opts.StateHandler.NodeAvailableMap()
	//	msg, _ := json.Marshal(nam)
	//	if _, err := c.Writer.Write(msg); err != nil {
	//		break
	//	}
	//
	//	// TODO should be configurable
	//	time.Sleep(3 * time.Second)
	//}

	c.Status(http.StatusOK)
}
