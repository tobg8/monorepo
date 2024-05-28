package logging

// NewNoop creates and returns a no-op logger.
func NewNoop() LoggerLevel {
	return &noop{LevelDebug}
}

type noop struct {
	lvl Level
}

func (n *noop) SetLevel(lvl Level) {
	n.lvl = lvl
}

func (n *noop) GetLevel() Level {
	return n.lvl
}

func (*noop) Debugf(format string, v ...interface{})   {}
func (*noop) Infof(format string, v ...interface{})    {}
func (*noop) Warningf(format string, v ...interface{}) {}
func (*noop) Errorf(format string, v ...interface{})   {}

func (*noop) Debug(v ...interface{})   {}
func (*noop) Info(v ...interface{})    {}
func (*noop) Warning(v ...interface{}) {}
func (*noop) Error(v ...interface{})   {}

func (n *noop) WithField(fn string, fv interface{}) Logger { return n }
func (n *noop) WithFields(newFields Fields) Logger         { return n }

func (n *noop) WithError(err error) Logger { return n }
