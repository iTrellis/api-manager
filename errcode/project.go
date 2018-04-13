package errcode

import "github.com/go-trellis/errors"

const namespaceProject = "api-manager::errcode::project"

var (
	ErrParseProjectRequest = errors.TN(namespaceProject, 10000, "parse request failed: {{.err}}")
)

var (
	ErrCreateProject = errors.TN(namespaceProject, 11001, "failed create project: {{.err}}")
	ErrUpdateProject = errors.TN(namespaceProject, 11002, "failed update project: {{.err}}")
)
