package main

import (
	"go.uber.org/zap"
)

type webserverLogger struct {
	*zap.Logger
}

func (wl *webserverLogger) expand(fields []zap.Field) []zap.Field {
	if len(gitBranch) > 0 {
		fields = append(fields, zap.String("git_branch", gitBranch))
	}
	if len(gitCommit) > 0 {
		fields = append(fields, zap.String("git_commit", gitCommit))
	}
	return fields
}

func (wl *webserverLogger) Debug(msg string, fields ...zap.Field) {
	wl.Logger.Info(msg, wl.expand(fields)...)
}

func (wl *webserverLogger) Info(msg string, fields ...zap.Field) {
	wl.Logger.Info(msg, wl.expand(fields)...)
}

func (wl *webserverLogger) Warn(msg string, fields ...zap.Field) {
	wl.Logger.Info(msg, wl.expand(fields)...)
}

func (wl *webserverLogger) Fatal(msg string, fields ...zap.Field) {
	wl.Logger.Info(msg, wl.expand(fields)...)
}
