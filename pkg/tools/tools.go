//go:build tools
// +build tools

// This file exists to track tool dependencies. This is one of the recommended practices
// for handling tool dependencies in a Go module as outlined here:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	// Install for hot reloading server
	_ "github.com/codegangsta/gin"
	// Install for managing the database
	_ "github.com/gobuffalo/pop/v6/soda"
	// Install for autogenerating mocks
	_ "github.com/vektra/mockery/v2"
	// Replacement for go-junit-report
	_ "gotest.tools/gotestsum"
	// Install for go-swagger code generation
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
)
