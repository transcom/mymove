package atolinter

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ATOAnalyzer describes an analysis function and its options.
var ATOAnalyzer = &analysis.Analyzer{
	Name:     "atolint",
	Doc:      "Checks that disabling of gosec is accompanied by annotations",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

const disableNoSec = "#nosec"
const disableErrcheck = "nolint:errcheck"

var linters = []string{disableNoSec, disableErrcheck}
var lintersString = strings.Join(linters, "|")

var lintersRegex = regexp.MustCompile(fmt.Sprintf("(?P<linterDisabled>%v)", lintersString))

const validatorStatusLabel = "RA Validator Status:"

var validatorStatuses = map[string]bool{
	"RA ACCEPTED":         true,
	"RETURN TO DEVELOPER": true,
	"KNOWN ISSUE":         true,
	"MITIGATED":           true,
	"FALSE POSITIVE":      true,
	"BAD PRACTICE":        true,
}

// check if comment group has disabling of gosec in it but it doesn't have a specific rule it is disabling
func containsGosecDisableNoRule(comments []*ast.Comment) bool {
	for _, comment := range comments {
		noSecRegex := regexp.MustCompile(fmt.Sprintf("(?P<linter>%v) ?(?P<rule>G\\d{3})?", disableNoSec))

		match := noSecRegex.FindStringSubmatch(comment.Text)

		if match == nil {
			return false
		}

		if match[2] == "" {
			return true
		}
	}
	return false
}

func containsDisabledLinterWithoutAnnotation(comments []*ast.Comment) bool {
	for _, comment := range comments {
		if strings.Contains(comment.Text, validatorStatusLabel) {
			return false
		}
	}
	return true
}

func findDisabledLinter(comments []*ast.Comment) (bool, string) {
	for _, comment := range comments {
		match := lintersRegex.FindStringSubmatch(comment.Text)

		if match == nil {
			continue
		}

		return true, match[1]
	}

	return false, ""
}

func containsAnnotationNotApproved(comments []*ast.Comment) bool {
	for _, comment := range comments {
		if strings.Contains(comment.Text, validatorStatusLabel) {
			// example str: //RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
			individualCommentArr := strings.Split(comment.Text, ": ")
			// Has validator status label but no value ex. //RA Validator Status:
			if len(individualCommentArr) == 1 {
				return true
			}
			for index, str := range individualCommentArr {
				str = strings.Trim(str, " ")
				if index > 0 && !validatorStatuses[strings.ToUpper(str)] {
					return true
				}
			}
		}
	}
	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	// pass.ResultOf[inspect.Analyzer] will be set if we've added inspect.Analyzer to Requires.
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.CommentGroup)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		comments := node.(*ast.CommentGroup)
		linterDisabled, linter := findDisabledLinter(comments.List)

		if !linterDisabled {
			return
		}

		if linter == disableNoSec {
			containsDisablingGosecWithNoRule := containsGosecDisableNoRule(comments.List)

			if containsDisablingGosecWithNoRule {
				pass.Reportf(node.Pos(), "Please provide the rule that is being disabled")
				return
			}
		}

		containsDisabledLinterWithoutAnnotation := containsDisabledLinterWithoutAnnotation(comments.List)
		if containsDisabledLinterWithoutAnnotation {
			pass.Reportf(node.Pos(), "Disabling of linter must have an annotation associated with it. Please visit https://docs.google.com/document/d/1qiBNHlctSby0RZeaPzb-afVxAdA9vlrrQgce00zjDww/edit#heading=h.b2vss780hqfi")
			return
		}

		containsAnnotationNotApproved := containsAnnotationNotApproved(comments.List)
		if containsAnnotationNotApproved {
			pass.Reportf(node.Pos(), "Annotation needs approval from an ISSO")
			return
		}
	})

	return nil, nil
}
