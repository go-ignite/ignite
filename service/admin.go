package service

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/model"
)

// --- admin

func (s *Service) AdminLogin(c *gin.Context) {
	req := new(api.AdminLoginRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	if req.Username != s.opts.Config.AdminUsername || req.Password != s.opts.Config.AdminPassword {
		s.errJSON(c, http.StatusUnauthorized, nil)
		return
	}

	token, err := s.createToken(s.opts.Config.AdminUsername)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, api.AdminLoginResponse{Token: token})
}

// --- account

func (s *Service) GetAccountList(c *gin.Context) {
	req := new(api.UserListRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	users, total, err := s.opts.ModelHandler.GetUserList(req.Keyword, req.PageIndex, req.PageSize)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	resp := make([]*api.User, len(users))
	for _, user := range users {
		resp = append(resp, &api.User{
			ID:   user.ID,
			Name: user.Name,
		})
	}

	c.JSON(http.StatusOK, api.UserListResponse{
		List:      resp,
		Total:     total,
		PageIndex: req.PageIndex,
	})
}

func (s *Service) DestroyAccount(c *gin.Context) {
	userID := c.Param("id")
	user, err := s.opts.ModelHandler.GetUserByID(userID)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	// TODO clean up service containers

	if err := s.opts.ModelHandler.DestroyUser(user.ID); err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// --- invite code

func (s *Service) GetInviteCodeList(c *gin.Context) {
	req := new(api.InviteCodeListRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	inviteCodes, total, err := s.opts.ModelHandler.GetAvailableInviteCodeList(req.PageIndex, req.PageSize)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	resp := make([]*api.InviteCode, len(inviteCodes))
	for _, ic := range inviteCodes {
		resp = append(resp, &api.InviteCode{
			ID:        ic.ID,
			Code:      ic.Code,
			Limit:     ic.Limit,
			ExpiredAt: ic.ExpiredAt.Unix(),
		})
	}

	c.JSON(http.StatusOK, api.InviteCodeListResponse{
		List:      resp,
		Total:     total,
		PageIndex: req.PageIndex,
	})
}

func (s *Service) RemoveInviteCode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}
	if err := s.opts.ModelHandler.DeleteInviteCode(id); err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (s *Service) GenerateInviteCodes(c *gin.Context) {
	req := new(api.GenerateCodesRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	var codes []*model.InviteCode
	for i := 0; i < int(req.Amount); i++ {
		codes = append(codes, model.NewInviteCode(req.Limit, time.Unix(req.ExpiredAt, 0)))
	}
	if err := s.opts.ModelHandler.CreateInviteCodes(codes); err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// --- node

// 1 -> cannot connect to node
// 2 -> cannot init node
// 3 -> node name already exists
func (s *Service) AddNode(c *gin.Context) {
	req := new(api.AddNodeRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}
	if s.opts.ModelHandler.CheckIfNodeNameExist(req.Name) {
		s.errJSON(c, http.StatusBadRequest, model.ErrDuplicateNodeName, 3)
		return
	}

	// TODO ping and init node, sync with node

	node := model.NewNode(req.Name, req.Comment, req.RequestAddress, req.ConnectionAddress, req.PortFrom, req.PortTo)
	if err := s.opts.ModelHandler.CreateNode(node); err != nil {
		if err == model.ErrDuplicateNodeName {
			s.errJSON(c, http.StatusBadRequest, model.ErrDuplicateNodeName, 3)
		} else {
			s.errJSON(c, http.StatusInternalServerError, errors.Wrap(err, "create node error"))
		}
		return
	}

	c.JSON(http.StatusCreated, &api.Node{
		ID:                node.ID,
		Name:              node.Name,
		Comment:           node.Comment,
		RequestAddress:    node.RequestAddress,
		ConnectionAddress: node.ConnectionAddress,
		PortFrom:          node.PortFrom,
		PortTo:            node.PortTo,
	})
}

func (s *Service) GetAllNodes(c *gin.Context) {
	nodes, err := s.opts.ModelHandler.GetAllNodes()
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	resp := make([]*api.Node, len(nodes))
	for _, node := range nodes {
		resp = append(resp, &api.Node{
			ID:                node.ID,
			Name:              node.Name,
			Comment:           node.Comment,
			RequestAddress:    node.RequestAddress,
			ConnectionAddress: node.ConnectionAddress,
			PortFrom:          node.PortFrom,
			PortTo:            node.PortTo,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Service) DeleteNode(c *gin.Context) {
	id := c.Param("id")
	if err := s.opts.ModelHandler.DeleteNode(id); err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	// TODO cancel node monitor

	c.JSON(http.StatusNoContent, nil)
}

// 1 -> there are services outside the port range
func (s *Service) UpdateNode(c *gin.Context) {
	req := new(api.UpdateNodeRequest)
	if err := c.ShouldBind(req); err != nil {
		s.errJSON(c, http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")
	count, err := s.opts.ModelHandler.GetNodePortRangeServiceCount(id, req.PortFrom, req.PortTo)
	if err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}
	if count > 0 {
		s.errJSON(c, http.StatusBadRequest, fmt.Errorf("there are %d services outside the port range", count), 1)
		return
	}

	node := &model.Node{
		ID:                id,
		Name:              req.Name,
		Comment:           req.Comment,
		ConnectionAddress: req.ConnectionAddress,
		PortFrom:          req.PortFrom,
		PortTo:            req.PortTo,
	}
	if err := s.opts.ModelHandler.SaveNode(node); err != nil {
		s.errJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
