package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/ss"
	"github.com/go-ignite/ignite/utils"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

var (
	servers          = []string{"SS", "SSR"}
	ssMethods        = []string{"aes-256-cfb", "aes-128-gcm", "aes-192-gcm", "aes-256-gcm", "chacha20-ietf-poly1305"}
	ssrMethods       = []string{"aes-256-cfb", "aes-256-ctr", "chacha20", "chacha20-ietf"}
	serverMethodsMap = map[string]map[string]bool{}
)

func init() {
	ssMethodMap := map[string]bool{}
	for _, method := range ssMethods {
		ssMethodMap[method] = true
	}
	ssrMethodMap := map[string]bool{}
	for _, method := range ssrMethods {
		ssrMethodMap[method] = true
	}

	serverMethodsMap["SS"] = ssMethodMap
	serverMethodsMap["SSR"] = ssrMethodMap
}

// ServiceConfigHandler godoc
// @Summary get user info
// @Description get user info
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} models.ServiceConfig
// @Router /api/user/auth/service/config [get]
func (router *MainRouter) ServiceConfigHandler(c *gin.Context) {
	sc := models.ServiceConfig{
		SSMethods:  ssMethods,
		SSRMethods: ssrMethods,
		Servers:    servers,
	}
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success!",
		Data:    sc,
	})
}

// PanelIndexHandler godoc
// @Summary get user info
// @Description get user info
// @Produce json
// @Success 200 {object} models.UserInfo
// @Param Authorization header string true "Authentication header"
// @Failure 200 {string} json "{"success":false,"message":"error message"}"
// @Router /api/user/auth/info [get]
func (router *MainRouter) UserInfoHandler(c *gin.Context) {
	userID, _ := c.Get("id")
	logrus.WithField("userID", userID).Debug("get user info")

	user := new(db.User)
	exists, err := router.db.Id(userID).Get(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"err":    err,
		}).Error("get user info error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取用户信息失败！"))
		return
	}

	if !exists {
		//Service has been removed by admininistrator.
		c.JSON(http.StatusOK, models.NewErrorResp("用户已删除！"))
		return
	}

	uInfo := &models.UserInfo{
		Id:            user.Id,
		Host:          ss.Host,
		Username:      user.Username,
		Status:        user.Status,
		PackageUsed:   fmt.Sprintf("%.2f", user.PackageUsed),
		PackageLimit:  user.PackageLimit,
		PackageLeft:   fmt.Sprintf("%.2f", float32(user.PackageLimit)-user.PackageUsed),
		ServicePort:   user.ServicePort,
		ServicePwd:    user.ServicePwd,
		ServiceMethod: user.ServiceMethod,
		ServiceType:   user.ServiceType,
		Expired:       user.Expired.Format("2006-01-02"),
		ServiceURL:    utils.ServiceURL(user.ServiceType, config.C.Host.Address, user.ServicePort, user.ServiceMethod, user.ServicePwd),
	}
	if uInfo.ServiceMethod == "" {
		uInfo.ServiceMethod = "aes-256-cfb"
	}
	if uInfo.ServiceType == "" {
		uInfo.ServiceType = "SS"
	}

	if user.PackageLimit == 0 {
		uInfo.PackageLeftPercent = "0"
	} else {
		uInfo.PackageLeftPercent = fmt.Sprintf("%.2f", (float32(user.PackageLimit)-user.PackageUsed)/float32(user.PackageLimit)*100)
	}

	logrus.WithField("userID", userID).Info("get info successful")
	c.JSON(http.StatusOK, models.NewSuccessResp(uInfo, "获取用户信息成功！"))
}

// CreateServiceHandler godoc
// @Summary create service
// @Description create service
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Param method formData string true "method"
// @Param server-type formData string true "server-type"
// @Success 200 {object} models.ServiceResult
// @Failure 200 {string} json "{"success":false,"message":"error message"}"
// @Router /api/user/auth/service/create [post]
func (router *MainRouter) CreateServiceHandler(c *gin.Context) {
	userID, _ := c.Get("id")
	method := c.PostForm("method")
	serverType := c.PostForm("server-type")

	logrus.WithFields(logrus.Fields{
		"userID":     userID,
		"method":     method,
		"serverType": serverType,
	}).Debug("create service")

	methodMap, ok := serverMethodsMap[serverType]
	if !ok {
		c.JSON(http.StatusOK, models.NewErrorResp("服务类型配置错误！"))
		return
	}

	if !methodMap[method] {
		c.JSON(http.StatusOK, models.NewErrorResp("加密方法配置错误！"))
		return
	}

	user := new(db.User)
	exists, err := router.db.Id(userID).Get(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"err":    err,
		}).Error("get user error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取用户信息失败！"))
		return
	}

	if !exists {
		//Service has been removed by admininistrator.
		c.JSON(http.StatusOK, models.NewErrorResp("用户已删除！"))
		return
	}
	if user.ServiceId != "" {
		c.JSON(http.StatusOK, models.NewErrorResp("服务已创建！"))
		return
	}

	//Get all used ports.
	var usedPorts []int
	if err := router.db.Table("user").Cols("service_port").Find(&usedPorts); err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"err":    err,
		}).Error("get used ports error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取用户服务端口失败！"))
	}
	logrus.WithField("usedPorts", usedPorts).Debug()

	// 1. Create ss service
	port, err := utils.GetAvailablePort(&usedPorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("创建服务失败，未能获取可用端口！"))
		return
	}
	logrus.WithField("port", port).Debug()
	result, err := ss.CreateAndStartContainer(serverType, strings.ToLower(user.Username), method, "", port)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":     userID,
			"serverType": serverType,
			"method":     method,
			"port":       port,
			"err":        err,
		}).Error("create service error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("创建服务失败！"))
		return
	}
	logrus.WithFields(logrus.Fields{
		"userID":     userID,
		"serverType": serverType,
		"method":     method,
		"port":       result.Port,
		"password":   result.Password,
		"serviceID":  result.ID,
	}).Info("create service success")

	// 2. Update user info
	user.Status = 1
	user.ServiceId = result.ID
	user.ServicePort = result.Port
	user.ServicePwd = result.Password
	user.ServiceMethod = method
	user.ServiceType = serverType
	affected, err := router.db.Id(userID).Cols("status", "service_port", "service_pwd", "service_id", "service_method", "service_type").Update(user)
	if affected == 0 || err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"err":    err,
		}).Error("update user service info error")

		//Force remove created container
		if err := ss.RemoveContainer(result.ID); err != nil {
			logrus.WithFields(logrus.Fields{
				"userID":    userID,
				"serviceID": result.ID,
			}).Error("remove service error")
		}

		c.JSON(http.StatusInternalServerError, models.NewErrorResp("更新用户信息失败！"))
		return
	}

	result.PackageLimit = user.PackageLimit
	result.Host = ss.Host

	c.JSON(http.StatusOK, models.NewSuccessResp(result, "服务创建成功！"))
}
