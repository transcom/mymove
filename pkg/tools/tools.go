// +build tools

// This file exists to track tool dependencies. This is one of the recommended practices
// for handling tool dependencies in a Go module as outlined here:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	_ "github.com/codegangsta/gin"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/gobuffalo/pop/soda"
	_ "github.com/securego/gosec/cmd/gosec"
	_ "github.com/segmentio/chamber"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/callgraph"
	_ "golang.org/x/tools/cmd/goimports"
)
