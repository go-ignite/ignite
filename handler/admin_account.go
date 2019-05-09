package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/state"
)

func (ah *AdminHandler) GetAccountList(c *gin.Context) {
	req := new(api.UserListRequest)
	if err := c.ShouldBind(req); err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}

	users, total, err := models.GetUserList(req.Keyword, req.PageIndex, req.PageSize)
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get account list error"))
		return
	}
	var resp []*api.User
	for _, user := range users {
		u := new(api.User)
		if err := ah.copy(u, user); err != nil {
			ah.ErrJSON(c, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, u)
	}

	c.JSON(http.StatusOK, api.NewUserListResponse(resp, total, req.PageIndex))
}

// 1 -> a node is unavailable
func (ah *AdminHandler) DestroyAccount(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}

	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get user error"))
		return
	}
	if user == nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	services, err := user.GetServices()
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get user's services error"))
		return
	}
	nodeIDs := services.GetNodeIDs()
	loader := state.GetLoader()
	for _, nodeID := range nodeIDs {
		ns := loader.GetNodeStatus(nodeID)
		if ns != nil {
			continue
		}
		if !ns.Available() {
			ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrapf(err, "node is unavailable, nodeID: %s", nodeID))
			return
		}
	}
	ctx := context.Background()
	for _, service := range services {
		ns := loader.GetNodeStatus(service.NodeID)
		if _, err := ns.Client().RemoveService(ctx, &protos.RemoveServiceRequest{
			Token:     ah.getToken(c),
			ServiceId: service.ContainerID,
		}); err != nil {
			ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "remove service error"))
			return
		}
	}

	if err := user.Destroy(services.IDs()); err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "destroy user error"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
