package controllers

import (
	"fmt"
	"log"
	"net/http"

	"ignite/models"
	"ignite/ss"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (router *MainRouter) PanelIndexHandler(c *gin.Context) {
	userID, exists := c.Get("userId")

	if !exists {
		c.HTML(http.StatusOK, "panel.html", nil)
		return
	}

	user := new(models.User)
	router.db.Id(userID).Get(user)

	uInfo := &models.UserInfo{
		Id:           user.Id,
		Host:         ss.Host,
		Username:     user.Username,
		Status:       user.Status,
		PackageUsed:  fmt.Sprintf("%.2f", user.PackageUsed),
		PackageLimit: user.PackageLimit,
		PackageLeft:  fmt.Sprintf("%.2f", float32(user.PackageLimit)-user.PackageUsed),
		ServicePort:  user.ServicePort,
		ServicePwd:   user.ServicePwd,
		Expired: 	  user.Expired.Format("2006-01-02"),
	}
	if user.PackageLimit == 0 {
		uInfo.PackageLeftPercent = "0"
	} else {
		uInfo.PackageLeftPercent = fmt.Sprintf("%.2f", (float32(user.PackageLimit)-user.PackageUsed)/float32(user.PackageLimit)*100)
	}

	c.HTML(http.StatusOK, "panel.html", gin.H{
		"uInfo": uInfo,
	})
}

func (router *MainRouter) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userId")
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func (router *MainRouter) CreateServiceHandler(c *gin.Context) {
	//time.Sleep(time.Second * 3)
	userID, _ := c.Get("userId")

	user := new(models.User)
	router.db.Id(userID).Get(user)

	// 1. Create ss service
	result, err := ss.CreateAndStartContainer(user.Username)

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
	affected, err := router.db.Id(userID).Cols("status", "service_port", "service_pwd", "service_id").Update(user)

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
