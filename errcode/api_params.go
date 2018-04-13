package errcode

import "github.com/go-trellis/errors"

const namespaceAPIParams = "api-manager::errcode::api-params"

var (
	ErrParseAPIParamsRequest = errors.TN(namespaceAPIParams, 10000, "parse request failed: {{.err}}")
)

var (
	ErrGetAPIParams    = errors.TN(namespaceAPIParams, 10001, "get api params failed: {{.aid}}, {{.err}}")
	ErrCreateAPIParams = errors.TN(namespaceAPIParams, 10002, "create api params failed: {{.err}}")
	ErrDeleteAPIParams = errors.TN(namespaceAPIParams, 10003, "delete api params failed: {{.err}}")
)
