package main

import "log"

func myLog(format string, args ...interface{}) {
	// test
	// testtest
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}
