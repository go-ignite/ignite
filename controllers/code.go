package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/utils"
)

func (router *MainRouter) InviteCodeListHandler(c *gin.Context) {
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

	codes := new([]*db.InviteCode)
	router.db.Desc("created").Where("available = 1").Limit(pageSize, pageSize*(pageIndex-1)).Find(codes)

	code := new(db.InviteCode)
	total, _ := router.db.Where("available = 1").Count(code)

	pd := models.PageData{Total: total, PageSize: pageSize, PageIndex: pageIndex, Data: codes}
	resp := models.Response{Success: true, Message: "success", Data: pd}
	c.JSON(http.StatusOK, resp)
}

func (router *MainRouter) RemoveInviteCodeHandler(c *gin.Context) {
	cid, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		resp := models.Response{Success: false, Message: "邀请码ID参数不正确"}
		c.JSON(http.StatusOK, resp)
		return
	}

	code := new(db.InviteCode)
	_, err = router.db.Id(cid).Delete(code)

	if err != nil {
		resp := models.Response{Success: false, Message: "邀请码删除失败"}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := models.Response{Success: true, Message: "success"}
	c.JSON(http.StatusOK, resp)
}

func (router *MainRouter) GenerateInviteCodeHandler(c *gin.Context) {
	generateCodeEntity := struct {
		Amount    int `json:"amount"`
		Limit     int `json:"limit"`
		Available int `json:"available"`
	}{}

	if err := c.BindJSON(&generateCodeEntity); err != nil {
		resp := models.Response{Success: false, Message: "Request body error..."}
		c.JSON(http.StatusOK, &resp)
		return
	}
	if generateCodeEntity.Amount == 0 || generateCodeEntity.Limit == 0 || generateCodeEntity.Available == 0 {
		resp := models.Response{Success: false, Message: "Data invalid..."}
		c.JSON(http.StatusOK, &resp)
		return
	}
	codes := []db.InviteCode{}
	for i := 0; i < generateCodeEntity.Amount; i++ {
		codes = append(codes, db.InviteCode{
			InviteCode:     utils.NewPasswd(16),
			PackageLimit:   generateCodeEntity.Limit,
			AvailableLimit: generateCodeEntity.Available,
			Available:      true,
		})
	}
	resp := models.Response{}
	if _, err := router.db.Insert(&codes); err != nil {
		log.Println("Save code error: ", err.Error())
		resp.Message = "Save codes error..."
	} else {
		resp.Success = true
	}
	c.JSON(http.StatusOK, &resp)
}
