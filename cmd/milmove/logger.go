package main

import (
	"fmt"

	"go.uber.org/zap"
)

type logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

//GorillaLogger wrapper for zap logger to use with gorilla recovery middleware
type GorillaLogger struct {
	*zap.Logger
}

//Println implementation of interface required by gorilla recovery handler
func (gl *GorillaLogger) Println(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	gl.Error(msg)
}
