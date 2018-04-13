package errcode

import "github.com/go-trellis/errors"

const namespaceAPI = "api-manager::errcode::api"

var (
	ErrParseAPIRequest = errors.TN(namespaceAPI, 10000, "parse request failed: {{.err}}")
)

var (
	ErrCreateAPI = errors.TN(namespaceAPI, 10001, "create api failed: {{.err}}")
	ErrGetAPI    = errors.TN(namespaceAPI, 10002, "get api failed: {{.id}}, {{.err}}")
	ErrUpdateAPI = errors.TN(namespaceAPI, 10003, "update api failed: {{.err}}")
)
