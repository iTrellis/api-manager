package errcode

import "github.com/go-trellis/errors"

const namespaceProject = "api-manager::errcode::project"

var (
	ErrParseProjectRequest = errors.TN(namespaceProject, 10000, "parse request failed: {{.err}}")
)

var (
	ErrGetProject    = errors.TN(namespaceProject, 10001, "failed get project: {{.err}}")
	ErrCreateProject = errors.TN(namespaceProject, 10002, "failed create project: {{.err}}")
	ErrUpdateProject = errors.TN(namespaceProject, 10003, "failed update project: {{.err}}")
)
