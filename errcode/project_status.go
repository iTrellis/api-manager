package errcode

import "github.com/go-trellis/errors"

const namespaceProjectStatus = "api-manager::errcode::project_status"

var (
	ErrGetProjectStatus = errors.TN(namespaceProjectStatus, 11001, "failed get project status: {{.err}}")
)
