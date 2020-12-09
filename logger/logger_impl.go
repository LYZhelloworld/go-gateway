package logger

import (
	"fmt"
	"io"
	"time"
)

type loggerImpl struct {
	// lowestLevel is the lowest level that the logger outputs (including itself).
	// For example, setting it as INFO will make the logger contains only INFO, WARN, ERROR, FATAL levels.
	lowestLevel Level
	// w is the writer where log outputs.
	w io.Writer
	// field
	dataField field
	// prev is a pointer to the previous loggerImpl. Log messages will be generated recursively.
	prev *loggerImpl
}

// getDataFields generates data fields recursively.
func (l loggerImpl) getDataFields() string {
	var msg string
	if l.prev != nil {
		msg = l.prev.getDataFields()
	}
	if l.dataField != nil {
		msg += fmt.Sprintf("%s=%v|", l.dataField.getKey(), l.dataField.getValue())
	}
	return msg
}

func (l loggerImpl) Enabled(level Level) bool {
	return l.lowestLevel <= level
}

func (l loggerImpl) outputWithLevel(level string, msg string) {
	now := time.Now().Format("2006-01-02 15:04:05-0700")
	_, _ = l.w.Write([]byte(fmt.Sprintf("%s|%s|%s%s\n", now, level, l.getDataFields(), msg)))
}

func (l loggerImpl) Trace(msg string) {
	if l.Enabled(Trace) {
		l.outputWithLevel("TRACE", msg)
	}
}

func (l loggerImpl) Debug(msg string) {
	if l.Enabled(Debug) {
		l.outputWithLevel("DEBUG", msg)
	}
}

func (l loggerImpl) Info(msg string) {
	if l.Enabled(Info) {
		l.outputWithLevel("INFO", msg)
	}
}

func (l loggerImpl) Warn(msg string) {
	if l.Enabled(Warn) {
		l.outputWithLevel("WARN", msg)
	}
}

func (l loggerImpl) Error(msg string) {
	if l.Enabled(Error) {
		l.outputWithLevel("ERROR", msg)
	}
}

func (l loggerImpl) Fatal(msg string) {
	if l.Enabled(Fatal) {
		l.outputWithLevel("FATAL", msg)
	}
}

func (l loggerImpl) WithField(key string, value interface{}) Logger {
	return &loggerImpl{
		lowestLevel: l.lowestLevel,
		w:           l.w,
		dataField:   &normalField{key: key, value: value},
		prev:        &l,
	}
}

func (l loggerImpl) WithError(err error) Logger {
	return &loggerImpl{
		lowestLevel: l.lowestLevel,
		w:           l.w,
		dataField:   &errorField{err: err},
		prev:        &l,
	}
}

func (l loggerImpl) WithLazyField(key string, value func() interface{}) Logger {
	return &loggerImpl{
		lowestLevel: l.lowestLevel,
		w:           l.w,
		dataField:   &lazyField{key: key, value: value},
		prev:        &l,
	}
}
