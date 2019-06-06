package service

//func (ah *Service) RemoveService(c *gin.Context) {
//	removeService(c, ah.logger)
//}
//
//func (ah *Service) ListServices(c *gin.Context) {
//	var (
//		userID, nodeID int64
//		err            error
//	)
//	if userIDStr := c.Query("userID"); userIDStr != "" {
//		if userID, err = strconv.ParseInt(userIDStr, 10, 64); err != nil {
//			c.JSON(http.StatusBadRequest, models.NewErrorResp(err))
//			return
//		}
//	}
//	if nodeIDStr := c.Query("nodeID"); nodeIDStr != "" {
//		if nodeID, err = strconv.ParseInt(nodeIDStr, 10, 64); err != nil {
//			c.JSON(http.StatusBadRequest, models.NewErrorResp(err))
//			return
//		}
//	}
//
//	listServices(c, userID, nodeID, ah.logger)
//}
