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

var defaultHFieldType *HFieldType

// HFieldType 操作者
type HFieldType struct {
	Committer     connector.Committer  `inject:"t"`
	RepoFieldType models.RepoFieldType `inject:"t"`
}

// NewHFieldType 获取操作对象
func NewHFieldType() *HFieldType {
	if defaultHFieldType == nil {
		defaultHFieldType = &HFieldType{}
	}
	return defaultHFieldType
}

// Inject 初始化函数
func (p *HFieldType) Inject(params map[string]interface{}) {
	injector.Inject(defaultHFieldType,
		params["committer"], params[models.ModelFieldTypeName])
}

/******************************
*** Method: Get
******************************/

type RespGetFieldType struct {
	CommonResponse
}

type RespGetFieldTypeItem struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (p *HFieldType) Get(ctx *gin.Context) {
	resp := RespGetFieldType{}

	var pts []domain.FieldType
	if err := p.Committer.NonTX(func(pRepo models.RepoFieldType) (ie error) {
		pts, ie = pRepo.GetList(nil)
		return
	}, p.RepoFieldType); err != nil {
		errResp := errcode.ErrGetFieldType.New(errors.Params{"err": err.Error()})
		resp.Code = errResp.Code()
		resp.Message = errResp.Error()
		resp.Namespace = errResp.Namespace()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	data := []RespGetFieldTypeItem{}
	for _, v := range pts {
		data = append(data, RespGetFieldTypeItem{
			ID:          v.ID,
			Type:        v.Type,
			Description: v.Description,
		})
	}
	resp.Data = data
	ctx.JSON(http.StatusOK, resp)
}
