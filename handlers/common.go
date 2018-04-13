package handlers

import (
	"github.com/go-playground/validator"
)

var Validator = validator.New()

type CommonResponse struct {
	Code      uint64      `json:"code"`
	Namespace string      `json:"namespace,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

const (
	DefaultPageNumber = 20
)
