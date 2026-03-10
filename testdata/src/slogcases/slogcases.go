package slogcases

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

var (
	slogCtx    = context.Background()
	slogLogger = slog.New(slog.NewTextHandler(io.Discard, nil))
)

func makeMessage(suffix string) string {
	return "starting " + suffix
}

func validTopLevelCalls() {
	slog.Info("starting server")
	slog.Warn("cache miss")
	slog.Error("request failed")
	slog.Debug("retry scheduled")
}

func invalidTopLevelCalls() {
	slog.Info("Starting server")           // want `log message must start with a lowercase letter`
	slog.Warn("ошибка подключения")        // want `log message must be in English \(ASCII only\)`
	slog.Error("server started!")          // want `log message must not contain special characters or emoji`
	slog.Debug("password rotation failed") // want `log message may contain sensitive data`
}

func validLoggerMethodCalls() {
	slogLogger.Info("connection established")
	slogLogger.WarnContext(slogCtx, "worker stopped")
	slogLogger.Log(slogCtx, slog.LevelInfo, "query executed")
	slogLogger.LogAttrs(slogCtx, slog.LevelInfo, "request completed", slog.String("component", "api"))
}

func invalidLoggerMethodCalls() {
	slogLogger.Info("Connection established")                                                                // want `log message must start with a lowercase letter`
	slogLogger.WarnContext(slogCtx, "запрос завершен")                                                       // want `log message must be in English \(ASCII only\)`
	slogLogger.Log(slogCtx, slog.LevelWarn, "worker stopped?")                                               // want `log message must not contain special characters or emoji`
	slogLogger.LogAttrs(slogCtx, slog.LevelError, "auth secret rotated", slog.String("component", "worker")) // want `log message may contain sensitive data`
}

func validContextAndChainedCalls() {
	slog.InfoContext(slogCtx, "shutdown complete")
	slog.Log(slogCtx, slog.LevelInfo, "queue drained")
	slog.LogAttrs(slogCtx, slog.LevelInfo, "batch processed", slog.String("job", "sync"))
	slogLogger.With("component", "api").Info("request started")
	slogLogger.WithGroup("db").ErrorContext(slogCtx, "connection failed")
}

func invalidContextAndChainedCalls() {
	slog.InfoContext(slogCtx, "Shutdown complete")                                           // want `log message must start with a lowercase letter`
	slog.Log(slogCtx, slog.LevelInfo, "данные обновлены")                                    // want `log message must be in English \(ASCII only\)`
	slog.LogAttrs(slogCtx, slog.LevelInfo, "batch processed...", slog.String("job", "sync")) // want `log message must not contain special characters or emoji`
	slogLogger.With("component", "api").Info("api_key rotated")                              // want `log message may contain sensitive data`
	slogLogger.WithGroup("db").ErrorContext(slogCtx, "Token refresh failed")                 // want `log message must start with a lowercase letter` `log message may contain sensitive data`
}

func skippedDynamicCalls() {
	message := "Starting server"
	secret := "hidden"

	slog.Info(message)
	slog.Info(fmt.Sprintf("starting %s", "server"))
	slog.Info(makeMessage("server"))
	slog.Info("password: " + secret)
	slogLogger.Info(message)
	slogLogger.With("component", "api").Info(fmt.Sprintf("request %s", "started"))
}
