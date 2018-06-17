package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/db/api"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/state"
	"github.com/go-ignite/ignite/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

func (uh *UserHandler) GetServiceConfig(c *gin.Context) {
	c.JSON(http.StatusOK, models.NewSuccessResp(agent.GetServiceConfigs()))
}

func (uh *UserHandler) ListServices(c *gin.Context) {
	userID := int64(c.GetFloat64("id"))
	dbAPI := api.NewAPI()
	services, err := dbAPI.GetServicesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResp("获取服务列表失败！"))
		return
	}
	servicesInfo := make([]*models.ServiceInfoResp, 0, len(services))
	for _, service := range services {
		sir := new(models.ServiceInfoResp)
		copier.Copy(sir, service)
		servicesInfo = append(servicesInfo, sir)
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(servicesInfo, "获取服务列表成功！"))
}

func (uh *UserHandler) CreateService(c *gin.Context) {
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
	uh.WithFields(logrus.Fields{
		"userID":   userID,
		"nodeID":   req.NodeID,
		"type":     req.Type,
		"method":   req.Method,
		"password": req.Password,
	}).Info("create service")

	exists, err := dbAPI.CheckServiceExists(userID, req.NodeID)
	if err != nil {
		uh.WithError(err).Error("check service error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("检查服务失败！"))
		return
	}
	if exists {
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
		port := int(resp.Port)
		ns.UsedPortMap[port] = true
		return port, nil
	}()
	if err != nil {
		ns.WithError(err).Error("get available port error")
		c.JSON(http.StatusOK, models.NewErrorResp("获取节点可用端口失败！"))
		return
	}
	uh.WithField("port", port).Info("agent available port")

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
		ns.WithError(err).Error("create service error")
		c.JSON(http.StatusOK, models.NewErrorResp("创建代理服务失败！"))
		return
	}

	uh.WithFields(logrus.Fields{
		"userID":    userID,
		"serviceID": agentResp.ServiceId,
	}).Info("create service success")

	service := &db.Service{
		ServiceID: agentResp.ServiceId,
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
				ServiceId: service.ServiceID,
			})
		}()
		uh.WithFields(logrus.Fields{
			"affected": affected,
			"error":    err,
		}).Error("create service error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("创建服务失败！"))
		return
	}

	resp := new(models.ServiceInfoResp)
	copier.Copy(resp, service)
	c.JSON(http.StatusOK, models.NewSuccessResp(resp, "创建服务成功！"))
}

func (uh *UserHandler) RemoveService(c *gin.Context) {
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
	nodeID, err := strconv.ParseInt(c.Param("nodeID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
		return
	}

	ns := state.GetLoader().GetNode(nodeID)
	if ns == nil {
		c.JSON(http.StatusOK, models.NewErrorResp("节点不存在！"))
		return
	}
	if !ns.Available() {
		c.JSON(http.StatusOK, models.NewErrorResp("节点暂不可用！"))
		return
	}
	uh.WithFields(logrus.Fields{
		"userID": userID,
		"nodeID": nodeID,
		"id":     id,
	}).Info("remove service")

	api := api.NewAPI()
	service, err := api.GetServiceByIDAndUserIDAndNodeID(id, userID, nodeID)
	if err != nil {
		uh.WithError(err).Error("get service error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取服务失败！"))
		return
	}
	if service.Id == 0 {
		c.JSON(http.StatusOK, models.NewErrorResp("服务不存在！"))
		return
	}

	if service.ServiceID != "" {
		if _, err := ns.RemoveService(context.Background(), &pb.RemoveServiceRequest{
			Token:     c.GetString("token"),
			ServiceId: service.ServiceID,
		}); err != nil {
			ns.WithFields(logrus.Fields{
				"error":     err,
				"serviceID": service.ServiceID,
			}).Error("remove service error")
			c.JSON(http.StatusOK, models.NewErrorResp("删除代理服务失败！"))
			return
		}
	}
	if _, err := api.RemoveServiceByID(id); err != nil {
		uh.WithError(err).Error("remove service error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("删除服务失败！"))
		return
	}
	go ns.RemovePortFromUsedMap(service.Port)
	c.JSON(http.StatusOK, models.NewSuccessResp(nil, "删除服务成功！"))
}
