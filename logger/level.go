package logger

// Level is the log level.
type Level int

const (
	// Trace level is the lowest log level used for tracing the flow of execution.
	Trace Level = iota

	// Debug level is used for debugging.
	Debug

	// Info level is used to highlight the progress of application, information and other useful data.
	Info

	// Warn level is for potentially harmful message but not considered an error usually.
	Warn

	// Error level means that some error has happened, and the application may stop or still continue running.
	Error

	// Fatal level is usually for those most severe errors that cause the application to abort immediately.
	Fatal
)
