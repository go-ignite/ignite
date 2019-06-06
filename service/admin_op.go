package service

//func (ah *Service) StopService(c *gin.Context) {
//	uid, err := strconv.Atoi(c.Param("id"))
//
//	if err != nil {
//		resp := models.Response{Success: false, Message: "用户ID参数不正确"}
//		c.JSON(http.StatusOK, resp)
//		return
//	}
//
//	user := new(db.Service)
//	db.GetDB().Id(uid).Get(user)
//
//	//1. Stop user's container
//	if ss.IsContainerRunning(user.ServiceId) {
//		err = ss.StopContainer(user.ServiceId)
//
//		if err != nil {
//			resp := models.Response{Success: false, Message: "停止服务失败"}
//			c.JSON(http.StatusOK, resp)
//			return
//		}
//
//		//2. Update Service status
//		user.Status = 2
//		db.GetDB().Id(uid).Cols("status").Update(user)
//	}
//
//	resp := models.Response{Success: true, Message: "success"}
//	c.JSON(http.StatusOK, resp)
//}
//
//func (ah *Service) StartServiceHandler(c *gin.Context) {
//	uid, err := strconv.Atoi(c.Param("id"))
//
//	if err != nil {
//		resp := models.Response{Success: false, Message: "用户ID参数不正确"}
//		c.JSON(http.StatusOK, resp)
//		return
//	}
//
//	user := new(db.Service)
//	db.GetDB().Id(uid).Get(user)
//
//	//1. Start user's container
//	if !ss.IsContainerRunning(user.ServiceId) {
//		err = ss.StartContainer(user.ServiceId)
//
//		if err != nil {
//			resp := models.Response{Success: false, Message: "启动服务失败"}
//			c.JSON(http.StatusOK, resp)
//			return
//		}
//
//		//2. Update Service status
//		user.Status = 1
//		db.GetDB().Id(uid).Cols("status").Update(user)
//	} else if user.Status == 2 {
//		user.Status = 1
//		db.GetDB().Id(uid).Cols("status").Update(user)
//	}
//
//	resp := models.Response{Success: true, Message: "success"}
//	c.JSON(http.StatusOK, resp)
//}
//
//func (ah *Service) RenewServiceHandler(c *gin.Context) {
//	uid, err := strconv.Atoi(c.Param("id"))
//
//	if err != nil {
//		resp := models.Response{Success: false, Message: "用户ID参数不正确"}
//		c.JSON(http.StatusOK, resp)
//		return
//	}
//
//	renewEntity := struct {
//		Expired int64 `json:"expired"`
//	}{}
//
//	if err := c.BindJSON(&renewEntity); err != nil || renewEntity.Expired == 0 {
//		resp := models.Response{Success: false, Message: "续费参数不正确!"}
//		c.JSON(http.StatusOK, &resp)
//		return
//	}
//
//	user := new(db.Service)
//	db.GetDB().Id(uid).Get(user)
//	if user.Id == 0 {
//		resp := models.Response{Success: false, Message: "获取用户失败!"}
//		c.JSON(http.StatusOK, &resp)
//		return
//	}
//
//	expired, now := time.Unix(renewEntity.Expired, 0), time.Now()
//	if expired.Before(now) {
//		resp := models.Response{Success: false, Message: "服务到期时间不能早于当前时间!"}
//		c.JSON(http.StatusOK, &resp)
//		return
//	}
//
//	if user.ServiceId != "" {
//		if user.Expired.Before(now) && (user.Status == 2 || !ss.IsContainerRunning(user.ServiceId)) {
//			if err := ss.StartContainer(user.ServiceId); err != nil {
//				resp := models.Response{Success: false, Message: "启动服务失败!"}
//				c.JSON(http.StatusOK, &resp)
//				return
//			}
//			user.Status = 1
//		}
//	}
//	user.Expired = expired
//
//	if _, err := db.GetDB().Id(uid).Cols("expired", "status").Update(user); err != nil {
//		resp := models.Response{Success: false, Message: "更新过期时间失败!"}
//		c.JSON(http.StatusOK, &resp)
//		return
//	}
//
//	updateMap := map[string]interface{}{
//		"id":      user.Id,
//		"expired": user.Expired.Format("2006-01-02"),
//		"status":  user.Status,
//	}
//	resp := models.Response{Success: true, Message: "success", Data: updateMap}
//	c.JSON(http.StatusOK, resp)
//}
