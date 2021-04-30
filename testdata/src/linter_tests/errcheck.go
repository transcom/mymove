package linter_tests

// nolint:errcheck // want "Disabling of linter must have an annotation associated with it. Please visit https://github.com/transcom/mymove/wiki/guide-to-static-analysis-annotations-for-disabled-Linters#guide-to-static-analysis-annotations-for-disabled-linters"
func errcheckShouldHaveAnnotation() {}

//RA Summary: [linter] - [linter type code] - [Linter summary] // want "Please add the truss-is3 team as reviewers for this PR and ping the ISSO in #static-code-review Slack. Add label ‘needs-is3-review’ to this PR. For more info see https://github.com/transcom/mymove/wiki/guide-to-static-analysis-security-workflow#guide-to-static-analysis-security-workflow"
//RA: <Why did the linter flag this line of code?>
//RA: <Why is this line of code valuable?>
//RA: <What mitigates the risk of negative impact?>
//RA Developer Status:  {RA Request, RA Accepted, POA&M Request, POA&M Accepted, Mitigated, Need Developer Fix, False Positive, Bad Practice}
//RA Validator Status:  {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
//RA Modified Severity: CAT III
// nolint:errcheck
func errcheckAnnotationNotApprovedTemplate() {}

//RA Summary: gosec - errcheck - Unchecked return value // want "Please add the truss-is3 team as reviewers for this PR and ping the ISSO in #static-code-review Slack. Add label ‘needs-is3-review’ to this PR. For more info see https://github.com/transcom/mymove/wiki/guide-to-static-analysis-security-workflow#guide-to-static-analysis-security-workflow"
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status:
//RA Validator Status:
//RA Modified Severity: N/A
// nolint:errcheck
func errcheckAnnotationNotApprovedEmpty() {}

//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to close a local server connection to ensure a unit test server is not left running indefinitely
//RA: Given the functions causing the lint errors are used to close a local server connection for testing purposes, it is not deemed a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Validator: jneuner@mitre.org
//RA Modified Severity: N/A
// nolint:errcheck
func errcheckAnnotationApproved() {}
