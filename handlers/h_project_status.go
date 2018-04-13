package handlers

import (
	"net/http"

	"github.com/go-trellis/api-manager/errcode"
	"github.com/go-trellis/api-manager/i12e/injector"
	"github.com/go-trellis/api-manager/models"
	"github.com/go-trellis/api-manager/models/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/connector"
	"github.com/go-trellis/errors"
)

/***********************
// Handler Information
***********************/

var defaultHProjectStatus *HProjectStatus

// HProjectStatus 操作者
type HProjectStatus struct {
	Committer         connector.Committer      `inject:"t"`
	RepoProjectStatus models.RepoProjectStatus `inject:"t"`
}

// NewHProjectStatus 获取操作对象
func NewHProjectStatus() *HProjectStatus {
	if defaultHProjectStatus == nil {
		defaultHProjectStatus = &HProjectStatus{}
	}
	return defaultHProjectStatus
}

// Inject 初始化函数
func (p *HProjectStatus) Inject(params map[string]interface{}) {
	injector.Inject(defaultHProjectStatus,
		params["committer"], params[models.ModelProjectStatusName])
}

/******************************
*** Method: Get
******************************/

type RespGetProjectStatus struct {
	CommonResponse
}

type RespGetProjectStatusItem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (p *HProjectStatus) Get(ctx *gin.Context) {
	resp := RespGetProjectStatus{}

	var pts []domain.ProjectStatus
	if err := p.Committer.NonTX(func(pRepo models.RepoProjectStatus) (ie error) {
		pts, ie = pRepo.GetList(nil)
		return
	}, p.RepoProjectStatus); err != nil {
		errResp := errcode.ErrGetProjectStatus.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	data := []RespGetProjectStatusItem{}
	for _, v := range pts {
		data = append(data, RespGetProjectStatusItem{
			ID:          v.ID,
			Description: v.Description,
		})
	}
	resp.Data = data
	ctx.JSON(http.StatusOK, resp)
}
