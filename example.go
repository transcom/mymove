package main

import "log"

func myLog(format string, args ...interface{}) {
	//RA Summary: gosec - G101 - Password Management: Hardcoded Password
	//RA: This line was flagged because it detected use of the word "token"
	//RA: This line is used to identify the name of the token. GorillaCSRFToken is the name of the base CSRF token.
	//RA: This variable does not store an application token.
	//RA Developer Status: Mitigated
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	// #nosec G1231
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}
