package main

import "log"

func myLog(format string, args ...interface{}) {
	// test
	// #nosec
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}
