package handlers

import (
	"net/http"

	"github.com/go-trellis/api-manager/errcode"
	"github.com/go-trellis/api-manager/i12e"
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

var defaultHProject *HProject

// HProject 操作者
type HProject struct {
	Committer   connector.Committer `inject:"t"`
	RepoProject models.RepoProject  `inject:"t"`
}

// NewHProject 获取操作对象
func NewHProject() *HProject {
	if defaultHProject == nil {
		defaultHProject = &HProject{}
	}
	return defaultHProject
}

// Inject 初始化函数
func (p *HProject) Inject(params map[string]interface{}) {
	injector.Inject(defaultHProject,
		params["committer"],
		params[models.ModelProjectName],
	)
}

/******************************
*** Method: Get
******************************/

// QueryGetProject
type QueryGetProject struct {
	ID int `form:"id" validate:"required"`
}

// RespGetProjectItem
type RespGetProjectItem struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	Host             string `json:"host"`
	ContactName      string `json:"contact_name"`
	ContactCellphone string `json:"contact_cellphone"`
	Brokerage        int64  `json:"brokerage"`
	Deposit          int64  `json:"deposit"`
	Refund           int64  `json:"refund"`
	Status           string `json:"status"`
}

type RespGetProject struct {
	CommonResponse
}

func (p *HProject) Get(ctx *gin.Context) {
	resp := RespGetProject{}

	req := QueryGetProject{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errResp := errcode.ErrParseProjectRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := Validator.Struct(&req); err != nil {
		errResp := errcode.ErrParseProjectRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var project *domain.Project
	if err := p.Committer.NonTX(func(pRepo models.RepoProject) (ie error) {
		project, ie = pRepo.Get(map[string]interface{}{"id": req.ID})
		return
	}, p.RepoProject); err != nil {
		errResp := errcode.ErrGetProject.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Data = &RespGetProjectItem{
		ID:               project.ID,
		Name:             project.Name,
		Address:          project.Address,
		Host:             project.Host,
		ContactName:      project.ContactName,
		ContactCellphone: project.ContactCellphone,
		Brokerage:        project.Brokerage,
		Deposit:          project.Deposit,
		Refund:           project.Refund,
		Status:           project.Status,
	}
	ctx.JSON(http.StatusOK, resp)
}

/******************************
*** Method: Post
******************************/

// ReqPostProject
type ReqPostProject struct {
	Name             string `json:"name" validate:"required"`
	Address          string `json:"address,omitempty"`
	Host             string `json:"host,omitempty"`
	ContactName      string `json:"contact_name"`
	ContactCellphone string `json:"contact_cellphone"`
	Brokerage        int64  `json:"brokerage"`
	Deposit          int64  `json:"deposit,omitempty"`
	Refund           int64  `json:"refund,omitempty"`
}

type RespPostProject struct {
	CommonResponse
}

func (p *HProject) Post(ctx *gin.Context) {
	resp := RespPostProject{}

	req := ReqPostProject{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errResp := errcode.ErrParseProjectRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := Validator.Struct(&req); err != nil {
		errResp := errcode.ErrParseProjectRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := p.Committer.NonTX(func(pRepo models.RepoProject) error {
		return pRepo.Add(&domain.Project{
			Name:             req.Name,
			Address:          req.Address,
			Host:             req.Host,
			ContactName:      req.ContactName,
			ContactCellphone: req.ContactCellphone,
			Brokerage:        req.Brokerage,
			Deposit:          req.Deposit,
			Refund:           req.Refund,
			Status:           i12e.StatusProjectNormal,
		})
	}, p.RepoProject); err != nil {
		errResp := errcode.ErrCreateProject.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

/******************************
*** Method: Put
******************************/

// ReqPutProject
type ReqPutProject struct {
	ID               int    `json:"id" validate:"required"`
	Address          string `json:"address,omitempty"`
	Host             string `json:"host,omitempty"`
	ContactName      string `json:"contact_name"`
	ContactCellphone string `json:"contact_cellphone"`
	Brokerage        int64  `json:"brokerage"`
	Deposit          int64  `json:"deposit,omitempty"`
	Refund           int64  `json:"refund,omitempty"`
	Status           string `json:"status" validate:"required"`
}

type RespPutProject struct {
	CommonResponse
}

func (p *HProject) Put(ctx *gin.Context) {
	resp := RespPutProject{}

	req := ReqPutProject{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errResp := errcode.ErrParseProjectRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := Validator.Struct(&req); err != nil {
		errResp := errcode.ErrParseProjectRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// TODO 校验Status的有效性
	if err := p.Committer.NonTX(func(pRepo models.RepoProject) error {
		return pRepo.Update(req.ID, map[string]interface{}{
			"address":           req.Address,
			"host":              req.Host,
			"contact_name":      req.ContactName,
			"contact_cellphone": req.ContactCellphone,
			"brokerage":         req.Brokerage,
			"deposit":           req.Deposit,
			"refund":            req.Refund,
			"status":            req.Status,
		})
	}, p.RepoProject); err != nil {
		errResp := errcode.ErrUpdateProject.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
