package mixedcases

import (
	"context"
	"io"
	"log/slog"

	"go.uber.org/zap"
)

var (
	mixedCtx    = context.Background()
	mixedSlog   = slog.New(slog.NewTextHandler(io.Discard, nil))
	mixedZap    = zap.NewNop()
	mixedSugar  = mixedZap.Sugar()
)

func validBoundaryCases() {
	slog.Info("oauth flow started")
	slog.Warn("tokenizer warmup complete")
	mixedSlog.Error("author cache miss")
	mixedZap.Info("secretary rotation complete")
	mixedSugar.Infow("apikeys cache synced", "component", "auth")
}

func invalidBoundaryCases() {
	slog.Info("auth flow started") // want `log message may contain sensitive data`
	mixedSlog.Warn("token refresh failed") // want `log message may contain sensitive data`
	mixedZap.Error("api_key rotation failed") // want `log message may contain sensitive data`
	mixedSugar.Infow("password reset issued", "component", "auth") // want `log message may contain sensitive data`
}

func skippedDynamicConcatenationCases() {
	secret := "hidden"
	status := "failed"

	slog.Info("password: " + secret)
	mixedSlog.Warn("token refresh " + status)
	mixedZap.Info("auth state " + status)
	mixedSugar.Infow("secret value "+secret, "component", "auth")
}

func multipleViolationCases() {
	slog.Info("Password leaked!") // want `log message must start with a lowercase letter` `log message must not contain special characters or emoji` `log message may contain sensitive data`
	mixedSlog.Log(mixedCtx, slog.LevelError, "Тoken leaked?") // want `log message must start with a lowercase letter` `log message must be in English \(ASCII only\)` `log message must not contain special characters or emoji`
	mixedZap.Warn("Secret rotation failed!") // want `log message must start with a lowercase letter` `log message must not contain special characters or emoji` `log message may contain sensitive data`
	mixedSugar.Errorw("Auth token leaked?", "component", "auth") // want `log message must start with a lowercase letter` `log message must not contain special characters or emoji` `log message may contain sensitive data`
}