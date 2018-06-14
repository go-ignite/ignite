package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/state"
	"github.com/go-ignite/ignite/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/db/api"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

// GetServiceConfigs godoc
// @Summary get service configs
// @Description get service configs
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} models.ServiceConfig
// @Router /api/user/auth/services/config [get]
func (uh *UserHandler) GetServiceConfig(c *gin.Context) {
	c.JSON(http.StatusOK, models.NewSuccessResp(agent.GetServiceConfigs()))
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
// @Router /api/user/auth/nodes/:id/services [post]
func (uh *UserHandler) CreateServiceHandler(c *gin.Context) {
	userID := int64(c.GetFloat64("id"))
	dbAPI := api.NewAPI()
	user, err := dbAPI.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResp("获取用户失败！"))
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusOK, models.NewErrorResp("用户已删除！"))
		return
	}

	req := &models.CreateServiceReq{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		req.Password = utils.RandString(10)
	}

	if req.NodeID, err = strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
		return
	}
	uh.WithFields(logrus.Fields{
		"userID":   userID,
		"nodeID":   req.NodeID,
		"type":     req.Type,
		"method":   req.Method,
		"password": req.Password,
	}).Debug("create service")

	typeProto := pb.ServiceType_NOT_SET
	serviceConfigs := agent.GetServiceConfigs()
	for _, serviceConfig := range serviceConfigs {
		if serviceConfig.Type == req.Type {
			typeProto = serviceConfig.TypeProto
			findMethod := false
			for _, method := range serviceConfig.Methods {
				if method == req.Method {
					findMethod = true
				}
			}
			if !findMethod {
				c.JSON(http.StatusOK, models.NewErrorResp("服务加密方法错误！"))
				return
			}
			break
		}
	}
	if typeProto == pb.ServiceType_NOT_SET {
		c.JSON(http.StatusOK, models.NewErrorResp("服务类型错误！"))
		return
	}
	if exists, err := dbAPI.CheckServiceExists(userID, req.NodeID); err != nil || exists {
		uh.WithFields(logrus.Fields{
			"userID": userID,
			"nodeID": req.NodeID,
			"exists": exists,
			"error":  err,
		}).Error("service has been created")
		c.JSON(http.StatusOK, models.NewErrorResp("重复创建服务！"))
		return
	}

	ns := state.GetLoader().GetNode(req.NodeID)
	if ns == nil {
		c.JSON(http.StatusOK, models.NewErrorResp("节点不存在！"))
		return
	}
	if !ns.Available() {
		c.JSON(http.StatusOK, models.NewErrorResp("节点暂不可用！"))
		return
	}

	// get available port from agent
	token := c.GetString("token")
	port, err := func() (int, error) {
		ns.Lock()
		defer ns.Unlock()
		req := &pb.GetAvailablePortRequest{
			Token:     token,
			UsedPorts: ns.GetUsedPorts(),
			PortFrom:  int32(ns.Node.PortFrom),
			PortTo:    int32(ns.Node.PortTo),
		}
		resp, err := ns.GetAvailablePort(context.Background(), req)
		if err != nil {
			return 0, err
		}
		ns.UsedPortMap[int(resp.Port)] = true
		return int(resp.Port), nil
	}()
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResp("获取节点可用端口失败！"))
		return
	}
	logrus.WithField("port", port).Debug("agent available port")

	// create service from agent
	agentResp, err := ns.CreateService(context.Background(), &pb.CreateServiceRequest{
		Token:    token,
		Port:     int32(port),
		Type:     typeProto,
		Method:   req.Method,
		Password: req.Password,
		Name:     user.Username,
	})
	if err != nil {
		go ns.RemovePortFromUsedMap(port)
		logrus.WithError(err).Error("create service error")
		c.JSON(http.StatusOK, models.NewErrorResp("创建服务失败！"))
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID":    userID,
		"serviceID": agentResp.ServiceId,
	}).Info("create service success")

	service := &db.Service{
		ServiceId: agentResp.ServiceId,
		UserId:    userID,
		NodeId:    ns.Node.Id,
		Type:      int(typeProto),
		Port:      int(port),
		Password:  req.Password,
		Method:    req.Method,
		Status:    1, // TODO change to enum
	}
	if affected, err := dbAPI.CreateService(service); err != nil || affected == 0 {
		go func() {
			ns.RemovePortFromUsedMap(port)
			ns.RemoveService(context.Background(), &pb.RemoveServiceRequest{
				Token:     token,
				ServiceId: service.ServiceId,
			})
		}()
		logrus.WithFields(logrus.Fields{
			"affected": affected,
			"error":    err,
		}).Error("create service error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("保存服务失败！"))
		return
	}

	resp := &models.CreateServiceResp{
		Port: int(port),
	}
	copier.Copy(resp, req)
	c.JSON(http.StatusOK, models.NewSuccessResp(resp, "创建服务成功！"))
}
