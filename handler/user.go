package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/ss"
	"github.com/go-ignite/ignite/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/db/api"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	*logrus.Logger
}

func NewUserHandler(l *logrus.Logger) *UserHandler {
	return &UserHandler{
		Logger: l,
	}
}

// LoginHandler godoc
// @Summary user login
// @Description user login
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"success":true,"message":"Success!","data":"Bearer xxxx"}"
// @Failure 200 {string} json "{"success":false,"message":"error message"}"
// @Router /api/user/login [post]
func (uh *UserHandler) LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")

	uh.WithFields(logrus.Fields{
		"username": username,
		"pwd":      pwd,
	}).Debug()

	user, err := api.NewAPI().GetUserByUsername(username)
	if err != nil {
		uh.WithFields(logrus.Fields{
			"username": username,
			"err":      err,
		}).Error("get user error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取用户信息失败！"))
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusOK, models.NewErrorResp("用户名不存在！"))
		return
	}

	if bcrypt.CompareHashAndPassword(user.HashedPwd, []byte(pwd)) != nil {
		c.JSON(http.StatusOK, models.NewErrorResp("用户名或密码错误！"))
		return
	}

	token, err := utils.CreateToken(config.C.App.Secret, user.Id)
	if err != nil {
		uh.WithFields(logrus.Fields{
			"userID": user.Id,
			"err":    err,
		}).Error("generate token error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取用户信息失败！"))
		return
	}

	uh.WithField("userID", user.Id).Info("login successful")
	c.JSON(http.StatusOK, models.NewSuccessResp(token))
}

// SignupHandler godoc
// @Summary user signup
// @Description user signup
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param invite-code formData string true "invite-code"
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Param confirm-password formData string true "confirm-password"
// @Success 200 {string} json "{"success":true,"message":"Success!","data":"Bearer xxxx"}"
// @Failure 200 {string} json "{"success":false,"message":"error message"}"
// @Router /api/user/signup [post]
func (uh *UserHandler) SignupHandler(c *gin.Context) {
	inviteCode := c.PostForm("invite-code")
	username := c.PostForm("username")
	pwd := c.PostForm("password")
	confirmPwd := c.PostForm("confirm-password")

	uh.WithFields(logrus.Fields{
		"inviteCode": inviteCode,
		"username":   username,
		"pwd":        pwd,
		"confirmPwd": confirmPwd,
	}).Debug()

	matched, _ := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9_.-]+$", username)
	if !matched {
		c.JSON(http.StatusOK, models.NewErrorResp("用户名不合规，请重新输入！"))
		return
	}

	if pwd != confirmPwd {
		c.JSON(http.StatusOK, models.NewErrorResp("密码不一致，请重新输入！"))
		return
	}

	iv, err := db.GetAvailableInviteCode(inviteCode)
	if err != nil {
		uh.WithFields(logrus.Fields{
			"inviteCode": inviteCode,
			"err":        err,
		}).Error("get invite code error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("查询邀请码记录失败！"))
		return
	}
	if iv.Id == 0 {
		c.JSON(http.StatusOK, models.NewErrorResp("邀请码不存在！"))
		return
	}
	uh.WithField("inviteCodeID", iv.Id).Debug()

	user := new(db.User)
	count, err := db.GetDB().Where("username = ?", username).Count(user)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResp("检查用户名失败！"))
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, models.NewErrorResp("用户名已存在！"))
		return
	}

	//Create user account
	trans := db.GetDB().NewSession()
	defer trans.Close()

	trans.Begin()

	//1.Create user account
	user.Username = username
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	user.HashedPwd = hashedPass
	user.InviteCode = iv.InviteCode
	user.PackageLimit = iv.PackageLimit
	user.Expired = time.Now().AddDate(0, iv.AvailableLimit, 0)

	affected, err := trans.Insert(user)
	if err != nil || affected == 0 {
		trans.Rollback()
		uh.WithFields(logrus.Fields{
			"err":      err,
			"affected": affected,
		}).Error("user insert error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("用户注册失败！"))
		return
	}

	//2.Set invite code as used status
	iv.Available = false
	affected, err = trans.ID(iv.Id).Cols("available").Update(iv)

	if err != nil || affected == 0 {
		trans.Rollback()
		uh.WithFields(logrus.Fields{
			"err":      err,
			"affected": affected,
		}).Error("update inviteCode status error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("更改邀请码可用状态失败！"))
		return
	}

	if err := trans.Commit(); err != nil {
		uh.WithField("err", err).Error("trans commit error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("用户注册失败！"))
		return
	}

	token, err := utils.CreateToken(config.C.App.Secret, user.Id)
	if err != nil {
		uh.WithFields(logrus.Fields{
			"userID": user.Id,
			"err":    err,
		}).Error("generate token error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("用户注册失败！"))
		return
	}

	uh.WithFields(logrus.Fields{
		"username":     username,
		"userID":       user.Id,
		"inviteCode":   inviteCode,
		"inviteCodeID": iv.Id,
	}).Info("register successful")
	c.JSON(http.StatusOK, models.NewSuccessResp(token))
}

// PanelIndexHandler godoc
// @Summary get user info
// @Description get user info
// @Produce json
// @Success 200 {object} models.UserInfo
// @Param Authorization header string true "Authentication header"
// @Failure 200 {string} json "{"success":false,"message":"error message"}"
// @Router /api/user/auth/info [get]
func (uh *UserHandler) UserInfoHandler(c *gin.Context) {
	userID, _ := c.Get("id")
	logrus.WithField("userID", userID).Debug("get user info")

	user := new(db.User)
	exists, err := db.GetDB().Id(userID).Get(user)
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
