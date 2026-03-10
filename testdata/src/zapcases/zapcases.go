package zapcases

import (
	"fmt"

	"go.uber.org/zap"
)

var (
	zapLogger   = zap.NewNop()
	sugarLogger = zapLogger.Sugar()
)

func buildMessage(suffix string) string {
	return "starting " + suffix
}

func validLoggerCalls() {
	zapLogger.Debug("worker started")
	zapLogger.Info("request completed", zap.String("component", "api"))
	zapLogger.Warn("cache miss")
	zapLogger.Error("connection failed")
	zapLogger.With(zap.String("component", "db")).Info("query executed")
}

func invalidLoggerCalls() {
	zapLogger.Info("Request completed") // want `log message must start with a lowercase letter`
	zapLogger.Warn("ошибка запроса") // want `log message must be in English \(ASCII only\)`
	zapLogger.Error("request failed!") // want `log message must not contain special characters or emoji`
	zapLogger.Debug("secret rotation failed") // want `log message may contain sensitive data`
	zapLogger.With(zap.String("component", "db")).Info("Token refresh failed") // want `log message must start with a lowercase letter` `log message may contain sensitive data`
}

func validSugaredCalls() {
	sugarLogger.Debugw("worker started", "component", "api")
	sugarLogger.Infow("request completed", "component", "api")
	sugarLogger.Warnw("cache miss", "component", "cache")
	sugarLogger.Errorw("connection failed", "component", "db")
	zapLogger.Sugar().Infow("job completed", "job", "sync")
}

func invalidSugaredCalls() {
	sugarLogger.Infow("Request completed", "component", "api") // want `log message must start with a lowercase letter`
	sugarLogger.Warnw("ошибка запроса", "component", "api") // want `log message must be in English \(ASCII only\)`
	sugarLogger.Errorw("request failed?", "component", "db") // want `log message must not contain special characters or emoji`
	sugarLogger.Debugw("auth token rotated", "component", "auth") // want `log message may contain sensitive data`
	zapLogger.Sugar().Infow("Api_key rotated", "component", "auth") // want `log message must start with a lowercase letter` `log message may contain sensitive data`
}

func skippedPrintStyleCalls() {
	message := "Request completed"
	secret := "hidden"

	sugarLogger.Info(message)
	sugarLogger.Infof("request %s", "completed")
	sugarLogger.Infoln("request completed")
	sugarLogger.Info("password:", secret)
	sugarLogger.Info(fmt.Sprintf("request %s", "completed"))
	sugarLogger.Info(buildMessage("server"))
	zapLogger.Sugar().Info(message)
}