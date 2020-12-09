package logger

type nopLoggerImpl struct {
}

func (n nopLoggerImpl) Enabled(_ Level) bool {
	return false
}

func (n nopLoggerImpl) Trace(_ string) {
	return
}

func (n nopLoggerImpl) Debug(_ string) {
	return
}

func (n nopLoggerImpl) Info(_ string) {
	return
}

func (n nopLoggerImpl) Warn(_ string) {
	return
}

func (n nopLoggerImpl) Error(_ string) {
	return
}

func (n nopLoggerImpl) Fatal(_ string) {
	return
}

func (n nopLoggerImpl) WithField(_ string, _ interface{}) Logger {
	return nopLoggerImpl{}
}

func (n nopLoggerImpl) WithError(_ error) Logger {
	return nopLoggerImpl{}
}

func (n nopLoggerImpl) WithLazyField(_ string, _ func() interface{}) Logger {
	return nopLoggerImpl{}
}
