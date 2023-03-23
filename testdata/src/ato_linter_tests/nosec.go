package ato_linter_tests

// #nosec // want "Please provide the rule that is being disabled"
func nosecShouldProvideRule() {}

// #nosec G101 // want "Disabling of linter must have an annotation associated with it. Please visit https://transcom.github.io/mymove-docs/docs/dev/contributing/code-analysis/Guide-to-Static-Analysis-Annotations-for-Disabled-Linters"
func nosecShouldHaveAnnotation() {}

//RA Summary: [linter] - [linter type code] - [Linter summary] // want "Please add the truss-is3 team as reviewers for this PR and ping the ISSO in #static-code-review Slack. Add label ‘needs-is3-review’ to this PR. For more info see https://transcom.github.io/mymove-docs/docs/dev/contributing/code-analysis/Guide-to-Static-Analysis-Security-Workflow"
//RA: <Why did the linter flag this line of code?>
//RA: <Why is this line of code valuable?>
//RA: <What mitigates the risk of negative impact?>
//RA Developer Status:  {RA Request, RA Accepted, POA&M Request, POA&M Accepted, Mitigated, Need Developer Fix, False Positive, Bad Practice}
//RA Validator Status:  {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
//RA Modified Severity: CAT III
// #nosec G402
func nosecAnnotationNotApprovedTemplate() {}

//RA Summary: gosec - G501 - Weak cryptographic hash  // want "Please add the truss-is3 team as reviewers for this PR and ping the ISSO in #static-code-review Slack. Add label ‘needs-is3-review’ to this PR. For more info see https://transcom.github.io/mymove-docs/docs/dev/contributing/code-analysis/Guide-to-Static-Analysis-Security-Workflow"
//RA: This line was flagged because of the use of MD5 hashing
//RA: This line of code hashes the AWS object to be able to verify data integrity
//RA: Purpose of this hash is to protect against environmental risks, it does not
//RA: hash any sensitive user provided information such as passwords.
//RA: AWS S3 API requires use of MD5 to validate data integrity.
//RA Developer Status:
//RA Validator Status:
//RA Modified Severity: CAT III
// #nosec G501
func nosecAnnotationNotApprovedEmpty() {}

//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to close a local server connection to ensure a unit test server is not left running indefinitely
//RA: Given the functions causing the lint errors are used to close a local server connection for testing purposes, it is not deemed a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Validator: jneuner@mitre.org
//RA Modified Severity: N/A
// #nosec G307
func nosecAnnotationApproved() {}
