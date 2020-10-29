// +build tools

// This file exists to track tool dependencies. This is one of the recommended practices
// for handling tool dependencies in a Go module as outlined here:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	// Install for hot reloading server
	_ "github.com/codegangsta/gin"
	// Install for managing the database
	_ "github.com/gobuffalo/pop/v5/soda"

	// Install for autogenerating mocks
	_ "github.com/vektra/mockery/cmd/mockery"

	// Test packages
	_ "github.com/go-playground/locales"
	_ "github.com/go-playground/universal-translator"
	_ "github.com/imdario/mergo"
	_ "github.com/leodido/go-urn"
	_ "github.com/namsral/flag"
	_ "github.com/stretchr/objx"
	_ "github.com/tealeg/xlsx"
	_ "gopkg.in/go-playground/validator.v9"

	// Install go-junit-report for CirclCI test result report generation
	_ "github.com/jstemmer/go-junit-report"
)
