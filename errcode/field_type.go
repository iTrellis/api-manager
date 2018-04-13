package errcode

import "github.com/go-trellis/errors"

const namespaceFieldType = "api-manager::errcode::field_type"

var (
	ErrGetFieldType = errors.TN(namespaceFieldType, 11001, "failed get common field type: {{.err}}")
)
