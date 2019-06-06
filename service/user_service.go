package service

//func (uh *Service) GetServiceConfig(c *gin.Context) {
//	c.JSON(http.StatusOK, models.NewSuccessResp(agent.GetServiceConfigs()))
//}
//
//func (uh *Service) ListServices(c *gin.Context) {
//	userID := int64(c.GetFloat64("id"))
//	listServices(c, userID, 0, uh.logger)
//}
//
//func (uh *Service) CreateService(c *gin.Context) {
//	userID := int64(c.GetFloat64("id"))
//	dbAPI := api.NewAPI()
//
//	user, err := uh.verifyUser(dbAPI, userID)
//	if err != nil {
//		c.JSON(http.StatusOK, models.NewErrorResp(err))
//		return
//	}
//
//	req := &models.CreateServiceReq{}
//	if err := c.BindJSON(req); err != nil {
//		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
//		return
//	}
//	if strings.TrimSpace(req.Password) == "" {
//		req.Password = utils.RandString(10)
//	}
//
//	if req.NodeID, err = strconv.ParseInt(c.Param("nodeID"), 10, 64); err != nil {
//		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
//		return
//	}
//
//	typeProto := pb.ServiceType_NOT_SET
//	serviceConfigs := agent.GetServiceConfigs()
//	for _, serviceConfig := range serviceConfigs {
//		if serviceConfig.Type == req.Type {
//			typeProto = serviceConfig.Type
//			findMethod := false
//			for _, method := range serviceConfig.Methods {
//				if method == req.Method {
//					findMethod = true
//				}
//			}
//			if !findMethod {
//				c.JSON(http.StatusOK, models.NewErrorResp("服务加密方法错误！"))
//				return
//			}
//			break
//		}
//	}
//	if typeProto == pb.ServiceType_NOT_SET {
//		c.JSON(http.StatusOK, models.NewErrorResp("服务类型错误！"))
//		return
//	}
//	uh.logger.WithFields(logrus.Fields{
//		"userID":   userID,
//		"nodeID":   req.NodeID,
//		"type":     req.Type,
//		"method":   req.Method,
//		"password": req.Password,
//	}).Info("create Service")
//
//	exists, err := dbAPI.CheckServiceExists(userID, req.NodeID)
//	if err != nil {
//		uh.logger.WithError(err).Error("check Service error")
//		c.JSON(http.StatusInternalServerError, models.NewErrorResp("检查服务失败！"))
//		return
//	}
//	if exists {
//		c.JSON(http.StatusOK, models.NewErrorResp("重复创建服务！"))
//		return
//	}
//
//	ns, err := verifyNode(req.NodeID)
//	if err != nil {
//		c.JSON(http.StatusOK, models.NewErrorResp(err))
//		return
//	}
//
//	// get available port from agent
//	token := c.GetString("token")
//	port, err := func() (int, error) {
//		ns.Lock()
//		defer ns.Unlock()
//		req := &pb.GetAvailablePortRequest{
//			Token:     token,
//			UsedPorts: ns.GetUsedPorts(),
//			PortFrom:  int32(ns.Node.PortFrom),
//			PortTo:    int32(ns.Node.PortTo),
//		}
//		resp, err := ns.Client.GetAvailablePort(context.Background(), req)
//		if err != nil {
//			return 0, err
//		}
//		port := int(resp.Port)
//		ns.UsedPortMap[port] = true
//		return port, nil
//	}()
//	ns.Logger.Info(1111)
//	if err != nil {
//		ns.Logger.WithError(err).Error("get available port error")
//		c.JSON(http.StatusOK, models.NewErrorResp("获取节点可用端口失败！"))
//		return
//	}
//	uh.logger.WithField("port", port).Info("agent available port")
//
//	// create Service from agent
//	agentResp, err := ns.Client.CreateService(context.Background(), &pb.CreateServiceRequest{
//		Token:    token,
//		Port:     int32(port),
//		Type:     typeProto,
//		Method:   req.Method,
//		Password: req.Password,
//		Name:     user.Username,
//	})
//	if err != nil {
//		go ns.RemovePortFromUsedMap(port)
//		ns.Logger.WithError(err).Error("create Service error")
//		c.JSON(http.StatusOK, models.NewErrorResp("创建代理服务失败！"))
//		return
//	}
//
//	uh.logger.WithFields(logrus.Fields{
//		"userID":    userID,
//		"serviceID": agentResp.ServiceId,
//	}).Info("create Service success")
//
//	Service := &db.Service{
//		ServiceID: agentResp.ServiceId,
//		UserId:    userID,
//		NodeId:    ns.Node.Id,
//		Type:      int(typeProto),
//		Port:      int(port),
//		Password:  req.Password,
//		Method:    req.Method,
//		Status:    1, // TODO change to enum
//	}
//	if affected, err := dbAPI.CreateService(Service); err != nil || affected == 0 {
//		go func() {
//			ns.RemovePortFromUsedMap(port)
//			ns.Client.RemoveService(context.Background(), &pb.RemoveServiceRequest{
//				Token:     token,
//				ServiceId: Service.ServiceID,
//			})
//		}()
//		uh.logger.WithFields(logrus.Fields{
//			"affected": affected,
//			"error":    err,
//		}).Error("create Service error")
//		c.JSON(http.StatusInternalServerError, models.NewErrorResp("创建服务失败！"))
//		return
//	}
//
//	resp := new(models.ServiceInfoResp)
//	copier.Copy(resp, Service)
//	c.JSON(http.StatusOK, models.NewSuccessResp(resp, "创建服务成功！"))
//}
//
//func (uh *Service) RemoveService(c *gin.Context) {
//	removeService(c, uh.logger)
//}
