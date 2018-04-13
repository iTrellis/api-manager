package handlers

import (
	"fmt"
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

var defaultHAPI *HAPI

// HAPI 操作者
type HAPI struct {
	Committer     connector.Committer  `inject:"t"`
	RepoAPI       models.RepoAPI       `inject:"t"`
	RepoAPIParams models.RepoAPIParams `inject:"t"`
	RepoProject   models.RepoProject   `inject:"t"`
}

// NewHAPI 获取操作对象
func NewHAPI() *HAPI {
	if defaultHAPI == nil {
		defaultHAPI = &HAPI{}
	}
	return defaultHAPI
}

// Inject 初始化函数
func (p *HAPI) Inject(params map[string]interface{}) {
	injector.Inject(defaultHAPI,
		params["committer"],
		params[models.ModelAPIName],
		params[models.ModelAPIParamsName],
		params[models.ModelProjectName],
	)
}

/******************************
*** Method: Get
******************************/

type QueryGetAPI struct {
	ID int `form:"id" validate:"required"`
}

type RespGetAPI struct {
	CommonResponse
}

type RespGetAPIItem struct {
	ID           int64  `json:"id"`
	ProjectID    int    `json:"project_id"`
	ProjectName  string `json:"project_name"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	HttpMethodID int    `json:"http_method_id"`
	HttpMethod   string `json:"http_method"`
	Description  string `json:"description"`
	Status       string `json:"status"`

	RespGetAPIParamsData RespGetAPIParamsData `json:"params"`
}

func (p *HAPI) Get(ctx *gin.Context) {
	resp := RespGetAPI{}

	req := QueryGetAPI{}
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

	var api *domain.API
	var params []domain.APIParameters
	if err := p.Committer.NonTX(func(
		apiRepo models.RepoAPI, apiParamsRepo models.RepoAPIParams) (ie error) {
		if api, ie = apiRepo.Get(map[string]interface{}{"id": req.ID}); ie != nil {
			return
		}
		params, ie = apiParamsRepo.GetList(map[string]interface{}{"api_id": req.ID})
		return
	}, p.RepoAPI, p.RepoAPIParams); err != nil {
		errResp := errcode.ErrGetAPI.New(errors.Params{"id": req.ID, "err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	// 第一层表示属于哪类的参数
	// 第二层表示数据哪个KEY的参数
	// 说明不同类型的请求参数可以重名
	customInfor := make(map[int][]*RespGetAPIParamsDataItem, 0)
	for _, v := range params {
		item := &RespGetAPIParamsItem{
			ID:            v.ID,
			ParentID:      v.ParentID,
			APIID:         v.APIID,
			APIParamsType: v.APIParamsType,
			Key:           v.Key,
			FieldType:     v.FieldType,
			IsList:        v.IsList,
			Required:      v.Required,
			Description:   v.Description,
			Sample:        v.Sample,
		}

		typeParams := customInfor[v.APIParamsType]
		if typeParams == nil || v.ParentID == 0 {
			typeParams = append(typeParams, &RespGetAPIParamsDataItem{Item: item})
			customInfor[v.APIParamsType] = typeParams
			continue
		}

		for _, dataItem := range typeParams {
			if dataItem.Rang(item) {
				break
			}
		}
		customInfor[v.APIParamsType] = typeParams
	}

	data := RespGetAPIParamsData{}
	for k, vs := range customInfor {
		for _, item := range vs {
			levelTreeOrder(item)
			switch k {
			case i12e.ParamTypeHeader:
				data.Headers = append(data.Headers, item.RangItems()...)
			case i12e.ParamTypeQuery:
				data.Queries = append(data.Queries, item.RangItems()...)
			case i12e.ParamTypeRequest:
				data.Requests = append(data.Requests, item.RangItems()...)
			case i12e.ParamTypeResponse:
				data.Responses = append(data.Responses, item.RangItems()...)
			}
		}
	}

	// ctx.JSON(http.StatusOK, resp)

	ctx.HTML(http.StatusOK, "api/index.tmpl",
		gin.H{"Title": "API参数列表", "Data": RespGetAPIItem{
			ID:           api.ID,
			ProjectID:    api.ProjectID,
			ProjectName:  api.ProjectName,
			Name:         api.Name,
			Path:         api.Path,
			HttpMethod:   api.HttpMethod,
			HttpMethodID: api.HttpMethodID,
			Description:  api.Description,
			Status:       api.Status,

			RespGetAPIParamsData: data,
		}})
}

/******************************
*** Method: Post
******************************/

type ReqPostAPI struct {
	ProjectID    int    `json:"project_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Path         string `json:"path" validate:"required"`
	HttpMethodID int    `json:"http_method_id" validate:"required"`
	Description  string `json:"description"`
}

type RespPostAPI struct {
	CommonResponse
}

func (p *HAPI) Post(ctx *gin.Context) {
	resp := RespPostAPI{}

	req := ReqPostAPI{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
	} else if 0 == len(i12e.MapAPIHTTPMethod[req.HttpMethodID]) {
		errResp := errcode.ErrParseAPIRequest.New(errors.Params{"err": "http method id not found"})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := p.Committer.NonTX(func(pRepo models.RepoProject, apiRepo models.RepoAPI) error {
		project, ie := pRepo.Get(map[string]interface{}{"id": req.ProjectID})
		if ie != nil {
			return ie
		}

		return apiRepo.Add(&domain.API{
			ProjectID:    project.ID,
			ProjectName:  project.Name,
			Name:         req.Name,
			Path:         req.Path,
			HttpMethodID: req.HttpMethodID,
			HttpMethod:   i12e.MapAPIHTTPMethod[req.HttpMethodID],
			Description:  req.Description,
			Status:       i12e.StatusAPINormal,
		})
	}, p.RepoProject, p.RepoAPI); err != nil {
		errResp := errcode.ErrCreateAPI.New(errors.Params{"err": err.Error()})
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

type ReqPutAPI struct {
	ID           int64  `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Path         string `json:"path" validate:"required"`
	HttpMethodID int    `json:"http_method_id" validate:"required"`
	Description  string `json:"description"`
}

type RespPutAPI struct {
	CommonResponse
}

func (p *HAPI) Put(ctx *gin.Context) {
	resp := RespPostAPI{}

	req := ReqPutAPI{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
	} else if 0 == len(i12e.MapAPIHTTPMethod[req.HttpMethodID]) {
		errResp := errcode.ErrParseAPIRequest.New(errors.Params{"err": "http method id not found"})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := p.Committer.NonTX(func(apiRepo models.RepoAPI) error {
		api, ie := apiRepo.Get(map[string]interface{}{"id": req.ID})
		if ie != nil {
			return ie
		} else if api.Status != i12e.StatusAPINormal {
			return fmt.Errorf("can't update,status ist normal: %d, %s", api.ID, api.Status)
		}

		return apiRepo.Update(api, map[string]interface{}{
			"name":           req.Name,
			"path":           req.Path,
			"http_method_id": req.HttpMethodID,
			"http_method":    i12e.MapAPIHTTPMethod[req.HttpMethodID],
			"description":    req.Description})
	}, p.RepoAPI); err != nil {
		errResp := errcode.ErrCreateAPI.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

/******************************
*** Method: Delete
******************************/

type QueryDeleteAPI struct {
	ID int64 `form:"id" validate:"required"`
}

type RespDeleteAPI struct {
	CommonResponse
}

func (p *HAPI) Delete(ctx *gin.Context) {
	resp := RespDeleteAPI{}

	req := QueryDeleteAPI{}
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

	if err := p.Committer.TX(func(apiRepo models.RepoAPI) error {
		api, ie := apiRepo.Get(map[string]interface{}{"id": req.ID})
		if ie != nil {
			return ie
		} else if api.Status == i12e.StatusAPIDeleted {
			return nil
		}

		return apiRepo.Update(api, map[string]interface{}{"status": i12e.StatusAPIDeleted})
	}, p.RepoAPI); err != nil {
		errResp := errcode.ErrUpdateAPI.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
