package handlers

import (
	"net/http"
	"time"

	"github.com/go-trellis/api-manager/errcode"
	"github.com/go-trellis/api-manager/i12e/injector"
	"github.com/go-trellis/api-manager/models"
	"github.com/go-trellis/api-manager/models/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/connector"
	"github.com/go-trellis/errors"
	"github.com/go-trellis/formats"
)

/***********************
// Handler Information
***********************/

var defaultHProjects *HProjects

// HProjects 操作者
type HProjects struct {
	Committer   connector.Committer `inject:"t"`
	RepoProject models.RepoProject  `inject:"t"`
}

// NewHProjects 获取操作对象
func NewHProjects() *HProjects {
	if defaultHProjects == nil {
		defaultHProjects = &HProjects{}
	}
	return defaultHProjects
}

// Inject 初始化函数
func (p *HProjects) Inject(params map[string]interface{}) {
	injector.Inject(defaultHProjects,
		params["committer"], params[models.ModelProjectName])
}

/******************************
*** Method: Get
******************************/

// QueryGetProjects 请求对象
type QueryGetProjects struct {
	Status string `form:"status"`
	Page   int    `form:"page"`
	Number int    `form:"num"`
}

type RespGetProjects struct {
	CommonResponse
}

type RespGetProjectsItem struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Address          string    `json:"address"`
	Host             string    `json:"host"`
	ContactName      string    `json:"contact_name"`
	ContactCellphone string    `json:"contact_cellphone"`
	Brokerage        int64     `json:"brokerage"`
	Deposit          int64     `json:"deposit"`
	Refund           int64     `json:"refund"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
	CreateTime       string    `json:"create_time"`
	UpdateTime       string    `json:"update_time"`
}

func (p *HProjects) Get(ctx *gin.Context) {
	resp := RespGetProjects{}

	req := QueryGetProjects{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errResp := errcode.ErrParseProjectRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if req.Number <= 0 {
		req.Number = DefaultPageNumber
	}
	offset := 0
	if req.Page > 1 {
		offset = (req.Page - 1) * req.Number
	}

	var pts []domain.Project

	if err := p.Committer.NonTX(func(pRepo models.RepoProject) (ie error) {
		params := map[string]interface{}{}
		if 0 != len(req.Status) {
			params["status"] = req.Status
		}
		pts, ie = pRepo.GetList(params, offset, req.Number)
		return
	}, p.RepoProject); err != nil {
		resp.Code = 2
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	data := []RespGetProjectsItem{}
	for _, v := range pts {
		data = append(data, RespGetProjectsItem{
			ID:               v.ID,
			Name:             v.Name,
			Address:          v.Address,
			Host:             v.Host,
			ContactName:      v.ContactName,
			ContactCellphone: v.ContactCellphone,
			Brokerage:        v.Brokerage,
			Deposit:          v.Deposit,
			Refund:           v.Refund,
			Status:           v.Status,
			CreateTime:       formats.FormatDateTime(v.CreatedAt),
			UpdateTime:       formats.FormatDateTime(v.UpdatedAt),
		})
	}
	// resp.Data = data
	// ctx.Render(http.StatusOK, render.JSON{})
	ctx.HTML(http.StatusOK, "projects/index.tmpl",
		gin.H{"AppSubURL": "", "Title": "项目列表", "Projects": data})
}
