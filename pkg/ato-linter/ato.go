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
const disableStaticcheck = "lint:ignore"
const disableStaticcheckFile = "lint:file-ignore"

var linters = []string{disableNoSec, disableErrcheck, disableStaticcheck, disableStaticcheckFile}
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

func checkForDisabledRule(linter string, comments []*ast.Comment) bool {
	var rulePattern string

	switch linter {
	case disableNoSec:
		rulePattern = "G\\d{3}"
	case disableStaticcheck, disableStaticcheckFile:
		rulePattern = "S[AT]?\\d{4}"
	}

	return containsDisableWithoutRule(linter, rulePattern, comments)
}

func containsDisableWithoutRule(linter string, rulePattern string, comments []*ast.Comment) bool {
	for _, comment := range comments {
		regex := regexp.MustCompile(fmt.Sprintf("(?P<linter>%v) ?(?P<rule>%s)?", linter, rulePattern))

		match := regex.FindStringSubmatch(comment.Text)

		if match == nil {
			continue
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
	inspectorInstance := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.File)(nil),
	}

	inspectorInstance.Preorder(nodeFilter, func(node ast.Node) {
		file := node.(*ast.File)
		comments := file.Comments

		for _, commentGroup := range comments {
			linterDisabled, linter := findDisabledLinter(commentGroup.List)

			if !linterDisabled {
				continue
			}

			switch linter {
			case disableNoSec, disableStaticcheck, disableStaticcheckFile:
				missingDisabledRule := checkForDisabledRule(linter, commentGroup.List)

				if missingDisabledRule {
					pass.Reportf(commentGroup.Pos(), "Please provide the rule that is being disabled")
					continue
				}
			}

			containsDisabledLinterWithoutAnnotation := containsDisabledLinterWithoutAnnotation(commentGroup.List)
			if containsDisabledLinterWithoutAnnotation {
				pass.Reportf(commentGroup.Pos(), "Disabling of linter must have an annotation associated with it. Please visit https://transcom.github.io/mymove-docs/docs/dev/contributing/code-analysis/Guide-to-Static-Analysis-Annotations-for-Disabled-Linters")
				continue
			}

			containsAnnotationNotApproved := containsAnnotationNotApproved(commentGroup.List)
			if containsAnnotationNotApproved {
				pass.Reportf(commentGroup.Pos(), "Please add the truss-is3 team as reviewers for this PR and ping the ISSO in #static-code-review Slack. Add label ‘needs-is3-review’ to this PR. For more info see https://transcom.github.io/mymove-docs/docs/dev/contributing/code-analysis/Guide-to-Static-Analysis-Security-Workflow")
				continue
			}
		}
	})

	return nil, nil
}
