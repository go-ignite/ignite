package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/ss"
	"github.com/go-ignite/ignite/utils"
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

func (router *MainRouter) PanelIndexHandler(c *gin.Context) {
	userID, exists := c.Get("userId")

	if !exists {
		c.HTML(http.StatusOK, "panel.html", nil)
		return
	}

	user := new(models.User)
	exists, _ = router.db.Id(userID).Get(user)

	if !exists {
		//Service has been removed by admininistrator.
		session := sessions.Default(c)
		session.Delete("userId")
		session.Save()

		c.Redirect(http.StatusFound, "/")
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
		ServiceURL:    utils.ServiceURL(user.ServiceType, utils.HOST_Address, user.ServicePort, user.ServiceMethod, user.ServicePwd),
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

	c.HTML(http.StatusOK, "panel.html", gin.H{
		"uInfo":       uInfo,
		"ss_methods":  ssMethods,
		"ssr_methods": ssrMethods,
		"servers":     servers,
	})
}

func (router *MainRouter) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userId")
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func (router *MainRouter) CreateServiceHandler(c *gin.Context) {
	userID, _ := c.Get("userId")
	method := c.PostForm("method")
	serverType := c.PostForm("server-type")

	fmt.Println("UserID", userID)
	fmt.Println("ServerType:", serverType)
	fmt.Println("Method:", method)

	methodMap, ok := serverMethodsMap[serverType]
	if !ok {
		resp := models.Response{Success: false, Message: "服务类型配置错误!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	if !methodMap[method] {
		resp := models.Response{Success: false, Message: "加密方法配置错误!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	user := new(models.User)
	router.db.Id(userID).Get(user)
	if user.ServiceId != "" {
		resp := models.Response{Success: false, Message: "服务已创建!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	//Get all used ports.
	var usedPorts []int
	router.db.Table("user").Cols("service_port").Find(&usedPorts)

	// 1. Create ss service
	port, err := utils.GetAvailablePort(&usedPorts)
	if err != nil {
		resp := models.Response{Success: false, Message: "创建服务失败,没有可用端口!"}
		c.JSON(http.StatusOK, resp)
		return
	}
	result, err := ss.CreateAndStartContainer(serverType, strings.ToLower(user.Username), method, "", port)
	if err != nil {
		log.Println("Create ss service error:", err.Error())
		resp := models.Response{Success: false, Message: "创建服务失败!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	// 2. Update user info
	user.Status = 1
	user.ServiceId = result.ID
	user.ServicePort = result.Port
	user.ServicePwd = result.Password
	user.ServiceMethod = method
	user.ServiceType = serverType
	affected, err := router.db.Id(userID).Cols("status", "service_port", "service_pwd", "service_id", "service_method", "service_type").Update(user)
	if affected == 0 || err != nil {
		if err != nil {
			log.Println("Update user info error:", err.Error())
		}

		//Force remove created container
		ss.RemoveContainer(result.ID)

		resp := models.Response{Success: false, Message: "更新用户信息失败!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	result.PackageLimit = user.PackageLimit
	result.Host = ss.Host
	resp := models.Response{Success: true, Message: "服务创建成功!", Data: result}

	c.JSON(http.StatusOK, resp)
}
