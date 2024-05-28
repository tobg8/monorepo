package logging

import (
	"fmt"
	"strings"
)

// Level type
type Level uint8

// These are the different logging levels.
const (
	// LevelError level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError Level = iota

	// LevelWarning level. Non-critical entries that deserve eyes.
	LevelWarning

	// LevelInfo level. General operational entries about what's going on inside
	// the application.
	LevelInfo

	// LevelDebug level. Usually only enabled when debugging. Very verbose
	// logging.
	LevelDebug
)

// Convert the Level to a string. E.g. LevelDebug becomes "debug".
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarning, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logging Level: %q", lvl)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*level = l

	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case LevelDebug:
		return []byte("debug"), nil
	case LevelInfo:
		return []byte("info"), nil
	case LevelWarning:
		return []byte("warning"), nil
	case LevelError:
		return []byte("error"), nil
	}

	return nil, fmt.Errorf("not a valid logging level %d", level)
}
