package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/db/api"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/state"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

func verifyUser(dbAPI *api.API, userID int64) (*db.User, error) {
	user, err := dbAPI.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, fmt.Errorf("用户已删除")
	}
	return user, nil
}

func verifyNode(nodeID int64) (*state.NodeStatus, error) {
	ns := state.GetLoader().GetNode(nodeID)
	if ns == nil {
		return nil, fmt.Errorf("节点不存在！")
	}
	if !ns.Available() {
		return nil, fmt.Errorf("节点暂不可用！")
	}
	return ns, nil
}

func removeService(c *gin.Context, logger *logrus.Logger) {
	dbAPI := api.NewAPI()
	userID := int64(c.GetFloat64("id"))

	if userID > 0 {
		if _, err := verifyUser(dbAPI, userID); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResp(err))
			return
		}
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp(err))
		return
	}

	service, err := dbAPI.GetServiceByIDAndUserID(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResp(err))
		return
	}
	if service.Id == 0 {
		c.JSON(http.StatusNotFound, models.NewErrorResp("服务不存在"))
		return
	}

	ns, err := verifyNode(service.NodeId)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResp(err))
		return
	}

	logger.WithFields(logrus.Fields{
		"userID": userID,
		"nodeID": service.NodeId,
		"id":     id,
	}).Info("remove service")

	if service.ServiceID != "" {
		if _, err := ns.Client.RemoveService(context.Background(), &pb.RemoveServiceRequest{
			Token:     c.GetString("token"),
			ServiceId: service.ServiceID,
		}); err != nil {
			ns.Logger.WithFields(logrus.Fields{
				"error":     err,
				"serviceID": service.ServiceID,
			}).Error("remove service error")
			c.JSON(http.StatusOK, models.NewErrorResp("删除代理服务失败！"))
			return
		}
	}
	if _, err := dbAPI.RemoveServiceByID(id); err != nil {
		logger.WithError(err).Error("remove service error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("删除服务失败！"))
		return
	}
	ns.RemovePortFromUsedMap(service.Port)
	c.JSON(http.StatusOK, models.NewSuccessResp(nil, "删除服务成功！"))
}

func listServices(c *gin.Context, userID, nodeID int64, logger *logrus.Logger) {
	dbAPI := api.NewAPI()
	services, err := dbAPI.GetServicesByUserIDAndNodeID(userID, nodeID)
	if err != nil {
		logger.WithError(err).Error("get service list error")
		c.JSON(http.StatusOK, models.NewErrorResp("获取服务列表失败！"))
		return
	}
	servicesInfo := make([]*models.ServiceInfoResp, 0, len(services))
	for _, service := range services {
		sir := new(models.ServiceInfoResp)
		copier.Copy(sir, service)
		sir.Created = service.Created.Unix()
		servicesInfo = append(servicesInfo, sir)
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(servicesInfo, "获取服务列表成功！"))
}
