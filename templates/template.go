package templates

import (
	"html/template"
	"strconv"

	"github.com/go-trellis/api-manager/i12e"
)

var TemplatesFuncMap = template.FuncMap{
	"APIStatusOperationName": APIStatusOperationName,
	"StatusName":             StatusName,
	"FundI2A":                FundI2A,
}

func StatusName(status string) string {
	if name, ok := i12e.MapStatusName[status]; ok {
		return name
	}
	return "未知状态"
}

func APIStatusOperationName(status string) string {
	switch status {
	case i12e.StatusAPINormal:
		return "删除"
	default:
		return ""
	}
}

func FundI2A(fund int64) string {
	sFund := strconv.Itoa(int(fund))
	lenFund := len(sFund)
	if lenFund <= 2 {
		return "0." + sFund
	}
	return sFund[0:lenFund-2] + "." + sFund[lenFund-2:]
}
