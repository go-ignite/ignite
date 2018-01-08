package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/ss"
)

var (
	methods    = []string{"aes-256-cfb", "aes-128-gcm", "aes-192-gcm", "aes-256-gcm", "chacha20-ietf-poly1305"}
	methodsMap = map[string]bool{
		"aes-256-cfb":            true,
		"aes-128-gcm":            true,
		"aes-192-gcm":            true,
		"aes-256-gcm":            true,
		"chacha20-ietf-poly1305": true,
	}
)

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
		Expired:       user.Expired.Format("2006-01-02"),
	}
	if uInfo.ServiceMethod == "" {
		uInfo.ServiceMethod = "aes-256-cfb"
	}

	if user.PackageLimit == 0 {
		uInfo.PackageLeftPercent = "0"
	} else {
		uInfo.PackageLeftPercent = fmt.Sprintf("%.2f", (float32(user.PackageLimit)-user.PackageUsed)/float32(user.PackageLimit)*100)
	}

	c.HTML(http.StatusOK, "panel.html", gin.H{
		"uInfo":   uInfo,
		"methods": methods,
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

	fmt.Println("UserID", userID)
	fmt.Println("Method:", method)

	if !methodsMap[method] {
		resp := models.Response{Success: false, Message: "Invalid method value!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	user := new(models.User)
	router.db.Id(userID).Get(user)
	if user.ServiceId != "" {
		resp := models.Response{Success: false, Message: "Service already created!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	//Get all used ports.
	var usedPorts []int
	router.db.Table("user").Cols("service_port").Find(&usedPorts)

	// 1. Create ss service
	result, err := ss.CreateAndStartContainer(user.Username, method, &usedPorts)
	if err != nil {
		log.Println("Create ss service error:", err.Error())
		resp := models.Response{Success: false, Message: "Create service error!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	// 2. Update user info
	user.Status = 1
	user.ServiceId = result.ID
	user.ServicePort = result.Port
	user.ServicePwd = result.Password
	user.ServiceMethod = method
	affected, err := router.db.Id(userID).Cols("status", "service_port", "service_pwd", "service_id", "service_method").Update(user)

	if affected == 0 || err != nil {
		if err != nil {
			log.Println("Update user info error:", err.Error())
		}

		//Force remove created container
		ss.RemoveContainer(result.ID)

		resp := models.Response{Success: false, Message: "Create service error!"}
		c.JSON(http.StatusOK, resp)
		return
	}

	result.PackageLimit = user.PackageLimit
	result.Host = ss.Host
	resp := models.Response{Success: true, Message: "OK!", Data: result}

	c.JSON(http.StatusOK, resp)
}
