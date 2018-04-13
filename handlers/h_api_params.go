package handlers

import (
	"container/list"
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

var defaultHAPIParams *HAPIParams

// HAPIParams 操作者
type HAPIParams struct {
	Committer     connector.Committer  `inject:"t"`
	RepoAPIParams models.RepoAPIParams `inject:"t"`
}

// NewHAPIParams 获取操作对象
func NewHAPIParams() *HAPIParams {
	if defaultHAPIParams == nil {
		defaultHAPIParams = &HAPIParams{}
	}
	return defaultHAPIParams
}

// Inject 初始化函数
func (p *HAPIParams) Inject(params map[string]interface{}) {
	injector.Inject(defaultHAPIParams,
		params["committer"],
		params[models.ModelAPIParamsName],
	)
}

/******************************
*** Method: Get
******************************/
type MyQueue struct {
	List *list.List
}

func (queue *MyQueue) push(elem interface{}) {
	queue.List.PushBack(elem)
}
func (queue *MyQueue) pop() interface{} {
	if elem := queue.List.Front(); elem != nil {
		queue.List.Remove(elem)
		return elem.Value
	}
	return nil
}
func levelTreeOrder(node *RespGetAPIParamsDataItem) {
	var nlast *RespGetAPIParamsDataItem
	last := node
	queue := MyQueue{List: list.New()}
	queue.push(node)

	for queue.List.Len() > 0 {
		node := queue.pop().(*RespGetAPIParamsDataItem)

		for _, elem := range node.SubParams {
			queue.push(elem)
			nlast = elem
		}

		if last == node {
			last = nlast
		}
	}
}

type RespGetAPIParamsData struct {
	Headers   []*RespGetAPIParamsItem `json:"headers"`
	Queries   []*RespGetAPIParamsItem `json:"queries"`
	Requests  []*RespGetAPIParamsItem `json:"requests"`
	Responses []*RespGetAPIParamsItem `json:"responses"`
}
type RespGetAPIParamsDataItem struct {
	Item *RespGetAPIParamsItem `json:"item"` // 自己

	SubParams []*RespGetAPIParamsDataItem `json:"sub_params"` // 子参数，自有 FieldType == object才可以
}

type RespGetAPIParamsItem struct {
	ID            int64  `json:"id"`
	ParentID      int64  `json:"parent_id"`
	APIID         int64  `json:"api_id"`
	APIParamsType int    `json:"api_params_type"`
	Key           string `json:"key"`
	FieldType     string `json:"field_type"`
	IsList        bool   `json:"is_list"`
	Required      bool   `json:"required"`
	Description   string `json:"description"`
	Sample        string `json:"sample"`
}

func (p *RespGetAPIParamsDataItem) Rang(item *RespGetAPIParamsItem) bool {
	if p.Item.ID == item.ParentID {
		p.SubParams = append(p.SubParams, &RespGetAPIParamsDataItem{Item: item})
		return true
	}
	for _, v := range p.SubParams {
		if v.Item.ID == item.ParentID {
			v.SubParams = append(v.SubParams, &RespGetAPIParamsDataItem{Item: item})
			return true
		}
	}

	return false
}

func (p *RespGetAPIParamsDataItem) RangItems() []*RespGetAPIParamsItem {
	var items []*RespGetAPIParamsItem
	items = append(items, p.Item)
	for _, v := range p.SubParams {
		if len(v.SubParams) != 0 {
			items = append(items, v.RangItems()...)
		} else {
			item := *v
			items = append(items, item.Item)
		}
	}
	return items
}

/******************************
*** Method: Post
******************************/

type ReqPostAPIParams struct {
	APIID         int64  `form:"aid" validate:"required"`
	ParentID      int64  `form:"parent_id"`
	APIParamsType int    `json:"api_params_type" validate:"required"`
	Key           string `json:"key" validate:"required"`
	FieldType     string `json:"field_type" validate:"required"`
	IsList        bool   `json:"is_list"`
	Required      bool   `json:"required"`
	Description   string `json:"description"`
	Sample        string `json:"sample"`
}

type RespPostAPIParams struct {
	CommonResponse
}

func (p *HAPIParams) Post(ctx *gin.Context) {
	resp := RespPostAPIParams{}

	req := ReqPostAPIParams{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errResp := errcode.ErrParseAPIParamsRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := Validator.Struct(&req); err != nil {
		errResp := errcode.ErrParseAPIParamsRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := p.Committer.NonTX(func(apiRepo models.RepoAPIParams) error {
		// TODO: Check Field type is common type
		if req.ParentID != 0 {
			// TODO 如果请求中的父节点ID不为0，则需要校验下对应的类型
		}
		return apiRepo.Add(&domain.APIParameters{
			APIID:         req.APIID,
			ParentID:      req.ParentID,
			APIParamsType: req.APIParamsType,
			Key:           req.Key,
			FieldType:     req.FieldType,
			IsList:        req.IsList,
			Required:      req.Required,
			Description:   req.Description,
			Sample:        req.Sample,
		})
	}, p.RepoAPIParams); err != nil {
		errResp := errcode.ErrCreateAPIParams.New(errors.Params{"err": err.Error()})
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

type QueryDeleteAPIParams struct {
	ID int64 `form:"id" validate:"required"`
}

type RespDeleteAPIParams struct {
	CommonResponse
}

func (p *HAPIParams) Delete(ctx *gin.Context) {
	resp := RespDeleteAPIParams{}

	req := QueryDeleteAPIParams{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errResp := errcode.ErrParseAPIParamsRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := Validator.Struct(&req); err != nil {
		errResp := errcode.ErrParseAPIParamsRequest.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := p.Committer.NonTX(func(apiRepo models.RepoAPIParams) error {
		return apiRepo.Delete(req.ID)
	}, p.RepoAPIParams); err != nil {
		errResp := errcode.ErrDeleteAPIParams.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
