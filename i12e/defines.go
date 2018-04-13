package i12e

import "net/http"

var MapStatusName = map[string]string{}

const (
	StatusProjectNormal     = "normal"     // 正常
	StatusProjectDoing      = "doing"      // 制作中
	StatusProjectDelivering = "delivering" // 交付中（制作完成）
	StatusProjectDelivered  = "delivered"  // 交付完成（待结账）
	StatusProjectDone       = "done"       // 正常结束（结账完成）
	StatusProjectUnfinished = "unfinished" // 烂尾项目（任何问题结束）
)

const (
	StatusAPINormal  = "normal"
	StatusAPIDeleted = "deleted"
)

func init() {
	MapStatusName = map[string]string{
		StatusProjectNormal:     "正常",
		StatusProjectDoing:      "开发中",
		StatusProjectDelivering: "交付中",
		StatusProjectDelivered:  "待结账",
		StatusProjectDone:       "已结束",
		StatusProjectUnfinished: "烂尾项目",
	}
}

const (
	ParamTypeHeader = iota + 1
	ParamTypeQuery
	ParamTypeRequest
	ParamTypeResponse
)

const (
	HTTPMethodGet = iota + 1
	HTTPMethodHead
	HTTPMethodPost
	HTTPMethodPut
	HTTPMethodPatch
	HTTPMethodDelete
	HTTPMethodConnect
	HTTPMethodOptions
	HTTPMethodTrace
)

var (
	MapAPIHTTPMethod = map[int]string{
		HTTPMethodGet:     http.MethodGet,
		HTTPMethodHead:    http.MethodHead,
		HTTPMethodPost:    http.MethodPost,
		HTTPMethodPut:     http.MethodPut,
		HTTPMethodPatch:   http.MethodPatch,
		HTTPMethodDelete:  http.MethodDelete,
		HTTPMethodConnect: http.MethodConnect,
		HTTPMethodOptions: http.MethodOptions,
		HTTPMethodTrace:   http.MethodTrace,
	}
)
