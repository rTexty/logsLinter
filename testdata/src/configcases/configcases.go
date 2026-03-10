package configcases

import "log/slog"

func example() {
	slog.Info("Starting server")
	slog.Info("credential rotated") // want `log message may contain sensitive data`
	slog.Info("request completed")
}
