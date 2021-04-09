package atolinter

import (
	"go/ast"
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
const validatorStatusLabel = "RA Validator Status:"

var validatorStatuses = map[string]bool{
	"RA Accepted":         true,
	"Return to Developer": true,
	"Known Issue":         true,
	"Mitigated":           true,
	"False Positive":      true,
	"Bad Practice":        true,
}

// check if comment group has disabling of gosec in it but it doesn't have a specific rule it is disabling
func containsGosecDisableNoReason(comments []*ast.Comment) bool {
	for _, comment := range comments {
		if strings.Contains(comment.Text, disableNoSec) {
			individualCommentArr := strings.Split(comment.Text, " ")
			for index, str := range individualCommentArr {
				if str == disableNoSec && index == len(individualCommentArr)-1 {
					return true
				}
			}
		}
	}
	return false
}

func containsGosecNoAnnotation(comments []*ast.Comment) bool {
	for _, comment := range comments {
		if strings.Contains(comment.Text, validatorStatusLabel) {
			return false
		}
	}
	return true
}

func containsNosec(comments []*ast.Comment) bool {
	for _, comment := range comments {
		if strings.Contains(comment.Text, disableNoSec) {
			individualCommentArr := strings.Split(comment.Text, " ")
			for _, str := range individualCommentArr {
				if str == disableNoSec {
					return true
				}
			}
		}
	}
	return false
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
				if index > 0 && !validatorStatuses[str] {
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
		commentsContainNosec := containsNosec(comments.List)

		if !commentsContainNosec {
			return
		}

		containsDisablingGosecWithNoReason := containsGosecDisableNoReason(comments.List)

		if containsDisablingGosecWithNoReason {
			pass.Reportf(node.Pos(), "Please provide the gosec rule that is being disabled")
			return
		}

		containsDisablingGosecNoAnnotation := containsGosecNoAnnotation(comments.List)
		if containsDisablingGosecNoAnnotation {
			pass.Reportf(node.Pos(), "Disabling of gosec must have an annotation associated with it. Please visit https://docs.google.com/document/d/1qiBNHlctSby0RZeaPzb-afVxAdA9vlrrQgce00zjDww/edit#heading=h.b2vss780hqfi")
			return
		}

		containsAnnotationNotApproved := containsAnnotationNotApproved(comments.List)
		if containsAnnotationNotApproved {
			pass.Reportf(node.Pos(), "Annotation needs approval from an ISSO")
			return
		}

		return
	})

	return nil, nil
}
