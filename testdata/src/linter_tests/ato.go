package linter_tests

// #nosec // want "Please visit https://github.com/transcom/mymove/wiki/guide-to-static-analysis-annotations-for-disabled-Linters#guide-to-static-analysis-annotations-for-disabled-linters"
func nosecShouldHaveAnnotation() {}

//RA Summary: [linter] - [linter type code] - [Linter summary] // want "Please add the truss-is3 team as reviewers for this PR and ping the ISSO in #static-code-review Slack. Add label ‘needs-is3-review’ to this PR. For more info see https://github.com/transcom/mymove/wiki/guide-to-static-analysis-security-workflow#guide-to-static-analysis-security-workflow"
//RA: <Why did the linter flag this line of code?>
//RA: <Why is this line of code valuable?>
//RA: <What mitigates the risk of negative impact?>
//RA Developer Status: {RA Request, RA Accepted, POA&M Request, POA&M Accepted, Mitigated, Need Developer Fix, False Positive, Bad Practice}
//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
//RA Validator: jneuner@mitre.org
//RA Modified Severity:
// #nosec G100
func nosecAnnotationNotApproved() {}
