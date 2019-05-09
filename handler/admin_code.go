package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/models"
)

func (ah *AdminHandler) GetInviteCodeList(c *gin.Context) {
	req := new(api.InviteCodeListRequest)
	if err := c.ShouldBind(req); err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}

	inviteCodes, total, err := models.GetAvailableInviteCodeList(req.PageIndex, req.PageSize)
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get invite code list error"))
		return
	}
	var ics []*api.InviteCode
	for _, inviteCode := range inviteCodes {
		ic := new(api.InviteCode)
		if err := ah.copy(ic, inviteCode); err != nil {
			ah.ErrJSON(c, http.StatusInternalServerError, err)
			return
		}
		ics = append(ics, ic)
	}

	c.JSON(http.StatusOK, api.NewInviteCodeListResponse(ics, total, req.PageIndex))
}

func (ah *AdminHandler) RemoveInviteCode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}
	if err := models.DeleteInviteCodeByID(id); err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "remove invite code error"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (ah *AdminHandler) GenerateInviteCodes(c *gin.Context) {
	req := new(api.GenerateCodesRequest)
	if err := c.ShouldBind(req); err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}

	var codes []*models.InviteCode
	for i := 0; i < int(req.Amount); i++ {
		codes = append(codes, models.NewInviteCode(req.Limit, req.ExpiredAt.Time()))
	}
	if err := models.CreateInviteCodes(codes); err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "generate invite codes error"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
