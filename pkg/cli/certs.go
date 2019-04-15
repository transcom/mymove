package cli

import (
	"fmt"
	"regexp"
)

// ParseCertificates takes a certificate and parses it into an slice of individual certificates
func ParseCertificates(str string) []string {

	certFormat := "-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----"

	// https://tools.ietf.org/html/rfc7468#section-2
	//	- https://stackoverflow.com/questions/20173472/does-go-regexps-any-charcter-match-newline
	re := regexp.MustCompile("(?s)([-]{5}BEGIN CERTIFICATE[-]{5})(\\s*)(.+?)(\\s*)([-]{5}END CERTIFICATE[-]{5})")
	matches := re.FindAllStringSubmatch(str, -1)

	certs := make([]string, 0, len(matches))
	for _, m := range matches {
		// each match will include a slice of strings starting with
		// (0) the full match, then
		// (1) "-----BEGIN CERTIFICATE-----",
		// (2) whitespace if any,
		// (3) base64-encoded certificate data,
		// (4) whitespace if any, and then
		// (5) -----END CERTIFICATE-----
		certs = append(certs, fmt.Sprintf(certFormat, m[3]))
	}
	return certs
}
