package main

import (
	"github.com/go-trellis/api-manager/control"
	"github.com/go-trellis/api-manager/i12e/builder"
)

// program defines
var (
	ProgramVersion  string
	CompilerVersion string
	Author          string
	BuildTime       string
)

func main() {
	builder.Show(ProgramVersion, CompilerVersion, BuildTime, Author)
	control.MainEntry()
}
