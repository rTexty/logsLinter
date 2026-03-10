package zap

type Field struct{}

type Logger struct{}

type SugaredLogger struct{}

func NewNop() *Logger {
	return &Logger{}
}

func String(string, string) Field {
	return Field{}
}

func (*Logger) Debug(string, ...Field) {}

func (*Logger) Info(string, ...Field) {}

func (*Logger) Warn(string, ...Field) {}

func (*Logger) Error(string, ...Field) {}

func (logger *Logger) With(...Field) *Logger {
	return logger
}

func (*Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{}
}

func (*SugaredLogger) Debugw(string, ...any) {}

func (*SugaredLogger) Infow(string, ...any) {}

func (*SugaredLogger) Warnw(string, ...any) {}

func (*SugaredLogger) Errorw(string, ...any) {}

func (*SugaredLogger) Debug(...any) {}

func (*SugaredLogger) Info(...any) {}

func (*SugaredLogger) Warn(...any) {}

func (*SugaredLogger) Error(...any) {}

func (*SugaredLogger) Debugf(string, ...any) {}

func (*SugaredLogger) Infof(string, ...any) {}

func (*SugaredLogger) Warnf(string, ...any) {}

func (*SugaredLogger) Errorf(string, ...any) {}

func (*SugaredLogger) Debugln(...any) {}

func (*SugaredLogger) Infoln(...any) {}

func (*SugaredLogger) Warnln(...any) {}

func (*SugaredLogger) Errorln(...any) {}
