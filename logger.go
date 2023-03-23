package gremgo

// Logger indicates the type of logging functionality that gremgo requires.
type Logger interface {
	// Log is a structured logging operation that will emit a log event from
	// the key value pairs given to it.
	Log(kvs ...any) error
}

// NopLogger is a no-operation logger. It does nothing.
type NopLogger struct{}

// Log does nothing for the no-operation logger.
func (l NopLogger) Log(kvs ...any) error {
	return nil
}
