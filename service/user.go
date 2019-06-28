package service

import (
	"net/http"
	"sort"

	"github.com/go-ignite/ignite-agent/protos"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/model"
	"github.com/go-ignite/ignite/state"
)

func (s *Service) UserLogin(c *gin.Context) {
	req := new(api.UserLoginRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	user, err := s.opts.ModelHandler.GetUserByName(req.Username)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		s.errJSON(c, http.StatusUnauthorized, nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPwd, []byte(req.Password)); err != nil {
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

func (s *Service) UserRegister(c *gin.Context) {
	req := new(api.UserRegisterRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	user := model.NewUser(req.Username, hashedPass)
	if err := s.opts.ModelHandler.CreateUser(user, req.InviteCode); err != nil {
		switch err {
		case model.ErrInviteCodeNotExistOrUnavailable:
			s.errJSON(c, http.StatusPreconditionFailed, err, 1)
		case model.ErrInviteCodeExpired:
			s.errJSON(c, http.StatusPreconditionFailed, err, 2)
		case model.ErrUserNameExists:
			s.errJSON(c, http.StatusPreconditionFailed, err, 3)
		default:
			s.errJSON(c, http.StatusInternalServerError, err)
		}

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

func (s *Service) CreateService(c *gin.Context) {
	userID := c.GetString("id")

	req := &api.CreateServiceRequest{}
	if err := c.BindJSON(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	sc := &model.ServiceConfig{
		EncryptionMethod: req.EncryptionMethod,
	}
	service := model.NewService(userID, req.NodeID, req.Type, sc)

	f := func() error {
		return s.opts.StateHandler.AddService(c.Request.Context(), service)
	}
	if err := s.opts.ModelHandler.CreateService(service, f); err != nil {
		switch err {
		case model.ErrUserDeleted:
			s.errJSON(c, http.StatusUnauthorized, nil)
		case model.ErrServiceExists:
			s.errJSON(c, http.StatusPreconditionFailed, err, 1)
		case state.ErrNodeNotExist:
			s.errJSON(c, http.StatusBadRequest, err)
		case state.ErrNodeUnavailable:
			s.errJSON(c, http.StatusPreconditionFailed, err, 2)
		default:
			s.errJSON(c, http.StatusInternalServerError, err)
		}
	}

	c.JSON(http.StatusOK, service.Output())
}

func (s *Service) GetServices(c *gin.Context) {
	c.JSON(http.StatusOK, s.opts.StateHandler.GetUserServices(c.GetString("id")))
}

func (s *Service) GetServiceOptions(c *gin.Context) {
	sos := make([]*api.ServiceOptions, 0, len(protos.ServiceType_Enum_name))
	for _, t := range []protos.ServiceType_Enum{protos.ServiceType_SS_LIBEV, protos.ServiceType_SSR} {
		so := &api.ServiceOptions{
			Type: t,
		}
		for k := range protos.ServiceEncryptionMethod_Enum_name {
			m := protos.ServiceEncryptionMethod_Enum(k)
			if m == protos.ServiceEncryptionMethod_NOT_SET {
				continue
			}

			if t.Suit(m) {
				so.EncryptionMethods = append(so.EncryptionMethods, m)
			}
		}
		sos = append(sos, so)
	}

	for _, so := range sos {
		sort.Slice(so.EncryptionMethods, func(i, j int) bool {
			return so.EncryptionMethods[i] < so.EncryptionMethods[j]
		})
	}

	c.JSON(http.StatusOK, sos)
}
