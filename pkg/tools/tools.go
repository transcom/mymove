// +build tools

// This file exists to track tool dependencies. This is one of the recommended practices
// for handling tool dependencies in a Go module as outlined here:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	// Install for hot reloading server
	_ "github.com/codegangsta/gin"
	// Install for generating swagger code
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	// Install for managing the database
	_ "github.com/gobuffalo/pop/soda"

	// Install for getting access to production secrets
	_ "github.com/segmentio/chamber"
	// Install for pre-commit go-lint
	_ "golang.org/x/lint/golint"
	// Install for pre-commit circleci testing
	_ "golang.org/x/tools/cmd/callgraph"
	// Install for pre-commit go-imports
	_ "golang.org/x/tools/cmd/goimports"
	// Install for pre-commit go-vet
	_ "golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow"
	// Install for linting project files & pre-commit
	_ "github.com/golangci/golangci-lint"

	// Packr isn't actually a tool dependency, but it's an indirect dependency that `go vet` and `go mod tidy` disagreed about.
	// Adding it here is a way to ensure that it isn't tidied up from go.mod
	_ "github.com/gobuffalo/packr"
)
