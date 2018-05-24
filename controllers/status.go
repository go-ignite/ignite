package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"
)

func (router *MainRouter) PanelStatusListHandler(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	kw := c.Query("keyword")

	users := new([]*db.User)
	router.db.Desc("created").Where("username like ?", "%"+kw+"%").Limit(pageSize, pageSize*(pageIndex-1)).Find(users)
	for _, user := range *users {
		if user.ServiceType == "" {
			if user.ServiceId != "" {
				user.ServiceType = "SS"
			} else {
				user.ServiceType = "N/A"
			}
		}
	}

	user := new(db.User)
	total, _ := router.db.Count(user)

	pd := models.PageData{Total: total, PageSize: pageSize, PageIndex: pageIndex, Data: users}
	resp := models.Response{Success: true, Message: "success", Data: pd}
	c.JSON(http.StatusOK, resp)
}
