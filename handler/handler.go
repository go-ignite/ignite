package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/logger"
)

type Handler struct {
	logger *logger.Logger
}

func New(l *logger.Logger) *Handler {
	return &Handler{
		logger: l,
	}
}

func (h *Handler) ErrJSON(c *gin.Context, statusCode int, err error, codes ...int) {
	code := statusCode
	if len(codes) > 0 {
		code = codes[0]
	}
	message := http.StatusText(statusCode)
	if err != nil {
		message = err.Error()
	}

	resp := api.NewErrResponse(code, message)
	h.logger.WithFields(logrus.Fields{
		"resp":       resp,
		"statusCode": statusCode,
	}).Error(c.Request.URL.String())

	c.JSON(code, resp)
}

func (h *Handler) createToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.C.App.Secret))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenStr, nil
}

func (h *Handler) getToken(c *gin.Context) string {
	return c.GetString("token")
}

func (h *Handler) isRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

func (h *Handler) copy(toValue, fromValue interface{}) error {
	return copier.Copy(toValue, fromValue)
}
