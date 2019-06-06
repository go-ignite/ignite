package service

//func verifyUser(dbAPI *api.API, userID int64) (*db.Service, error) {
//	user, err := dbAPI.GetUserByID(userID)
//	if err != nil {
//		return nil, err
//	}
//	if user.Id == 0 {
//		return nil, fmt.Errorf("用户已删除")
//	}
//	return user, nil
//}
//
//func verifyNode(nodeID int64) (*state.NodeStatus, error) {
//	ns := state.GetLoader().GetNode(nodeID)
//	if ns == nil {
//		return nil, fmt.Errorf("节点不存在！")
//	}
//	if !ns.Available() {
//		return nil, fmt.Errorf("节点暂不可用！")
//	}
//	return ns, nil
//}
//
//func removeService(c *gin.Context, logger *logrus.Logger) {
//	dbAPI := api.NewAPI()
//	userID := int64(c.GetFloat64("id"))
//
//	if userID > 0 {
//		if _, err := verifyUser(dbAPI, userID); err != nil {
//			c.JSON(http.StatusOK, models.NewErrorResp(err))
//			return
//		}
//	}
//
//	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, models.NewErrorResp(err))
//		return
//	}
//
//	Service, err := dbAPI.GetServiceByIDAndUserID(id, userID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, models.NewErrorResp(err))
//		return
//	}
//	if Service.Id == 0 {
//		c.JSON(http.StatusNotFound, models.NewErrorResp("服务不存在"))
//		return
//	}
//
//	ns, err := verifyNode(Service.NodeId)
//	if err != nil {
//		c.JSON(http.StatusOK, models.NewErrorResp(err))
//		return
//	}
//
//	logger.WithFields(logrus.Fields{
//		"userID": userID,
//		"nodeID": Service.NodeId,
//		"id":     id,
//	}).Info("remove Service")
//
//	if Service.ServiceID != "" {
//		if _, err := ns.Client.RemoveService(context.Background(), &pb.RemoveServiceRequest{
//			Token:     c.GetString("token"),
//			ServiceId: Service.ServiceID,
//		}); err != nil {
//			ns.Logger.WithFields(logrus.Fields{
//				"error":     err,
//				"serviceID": Service.ServiceID,
//			}).Error("remove Service error")
//			c.JSON(http.StatusOK, models.NewErrorResp("删除代理服务失败！"))
//			return
//		}
//	}
//	if _, err := dbAPI.RemoveServiceByID(id); err != nil {
//		logger.WithError(err).Error("remove Service error")
//		c.JSON(http.StatusInternalServerError, models.NewErrorResp("删除服务失败！"))
//		return
//	}
//	ns.RemovePortFromUsedMap(Service.Port)
//	c.JSON(http.StatusOK, models.NewSuccessResp(nil, "删除服务成功！"))
//}
//
//func listServices(c *gin.Context, userID, nodeID int64, logger *logrus.Logger) {
//	dbAPI := api.NewAPI()
//	services, err := dbAPI.GetServicesByUserIDAndNodeID(userID, nodeID)
//	if err != nil {
//		logger.WithError(err).Error("get Service list error")
//		c.JSON(http.StatusOK, models.NewErrorResp("获取服务列表失败！"))
//		return
//	}
//	servicesInfo := make([]*models.ServiceInfoResp, 0, len(services))
//	for _, Service := range services {
//		sir := new(models.ServiceInfoResp)
//		copier.Copy(sir, Service)
//		sir.Created = Service.Created.Unix()
//		servicesInfo = append(servicesInfo, sir)
//	}
//	c.JSON(http.StatusOK, models.NewSuccessResp(servicesInfo, "获取服务列表成功！"))
//}
