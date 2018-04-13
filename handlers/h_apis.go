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

var defaultHAPIs *HAPIs

// HAPIs 操作者
type HAPIs struct {
	Committer connector.Committer `inject:"t"`
	RepoAPI   models.RepoAPI      `inject:"t"`
}

// NewHAPIs 获取操作对象
func NewHAPIs() *HAPIs {
	if defaultHAPIs == nil {
		defaultHAPIs = &HAPIs{}
	}
	return defaultHAPIs
}

// Inject 初始化函数
func (p *HAPIs) Inject(params map[string]interface{}) {
	injector.Inject(defaultHAPIs,
		params["committer"],
		params[models.ModelAPIName],
	)
}

/******************************
*** Method: Get
******************************/

type QueryGetAPIs struct {
	ProjectID int `form:"pid" validate:"required"`
	Page      int `form:"page"`
	Number    int `form:"num"`
}

type RespGetAPIs struct {
	CommonResponse
}

type RespGetAPIsItem struct {
	ID           int64  `json:"id"`
	ProjectID    int    `json:"project_id"`
	ProjectName  string `json:"project_name"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	HttpMethodID int    `json:"http_method_id"`
	HttpMethod   string `json:"http_method"`
	Description  string `json:"description"`
	Status       string `json:"status"`
}

func (p *HAPIs) Get(ctx *gin.Context) {
	resp := RespGetAPIs{}

	req := QueryGetAPIs{}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		errResp := errcode.ErrParseAPIRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := Validator.Struct(&req); err != nil {
		errResp := errcode.ErrParseAPIRequest.New(errors.Params{"err": err.Error()})
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

	var apis []domain.API

	if err := p.Committer.NonTX(func(apiRepo models.RepoAPI) (ie error) {
		params := map[string]interface{}{
			"project_id": req.ProjectID,
			"status":     i12e.StatusAPINormal,
		}
		apis, ie = apiRepo.GetList(params, offset, req.Number)
		return
	}, p.RepoAPI); err != nil {
		resp.Code = 2
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	data := []RespGetAPIsItem{}
	projectName := ""
	for _, v := range apis {
		projectName = v.ProjectName
		data = append(data, RespGetAPIsItem{
			ID:           v.ID,
			ProjectID:    v.ProjectID,
			ProjectName:  v.ProjectName,
			Name:         v.Name,
			Path:         v.Path,
			HttpMethodID: v.HttpMethodID,
			HttpMethod:   v.HttpMethod,
			Description:  v.Description,
			Status:       v.Status,
		})
	}
	resp.Data = data
	// ctx.JSON(http.StatusOK, resp)
	ctx.HTML(http.StatusOK, "apis/index.tmpl",
		gin.H{"Title": "API列表", "APIs": data, "ProjectName": projectName})
}
