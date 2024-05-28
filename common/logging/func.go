package logging

import (
	"fmt"
)

type funcLogger struct {
	lvl     Level
	logFunc func(v ...interface{})
	helper  func()
	fields  map[string]interface{}
}

// FromFunc creates a logger with the given log function.
func FromFunc(logFunc func(v ...interface{})) LoggerLevel {
	return &funcLogger{
		lvl:     LevelInfo,
		logFunc: logFunc,
		helper:  func() {},
	}
}

// FromFuncHelper creates a logger with the given log function.
func FromFuncHelper(logFunc func(v ...interface{}), helper func()) LoggerLevel {
	return &funcLogger{
		lvl:     LevelInfo,
		logFunc: logFunc,
		helper:  helper,
	}
}

func (fl *funcLogger) SetLevel(lvl Level) {
	fl.lvl = lvl
}

func (fl *funcLogger) GetLevel() Level {
	return fl.lvl
}

func (fl *funcLogger) Debugf(format string, v ...interface{}) {
	fl.helper()

	if fl.lvl >= LevelDebug {
		fl.logf("[debug] "+format, v...)
	}
}

func (fl *funcLogger) Infof(format string, v ...interface{}) {
	fl.helper()

	if fl.lvl >= LevelInfo {
		fl.logf("[info] "+format, v...)
	}
}

func (fl *funcLogger) Warningf(format string, v ...interface{}) {
	fl.helper()

	if fl.lvl >= LevelWarning {
		fl.logf("[warning] "+format, v...)
	}
}

func (fl *funcLogger) Errorf(format string, v ...interface{}) {
	fl.helper()
	fl.logf("[error] "+format, v...)
}

func (fl *funcLogger) logf(format string, v ...interface{}) {
	fl.helper()

	var formatFields string
	for fn, fv := range fl.fields {
		if isEmptyFieldValue(fv) {
			continue
		}

		if formatFields != "" {
			formatFields += " X "
		}
		formatFields += "{%s: %v}"

		v = append(v, fn, fv)
	}

	fl.logFunc(fmt.Sprintf(format+formatFields, v...))
}

func (fl *funcLogger) Debug(v ...interface{}) {
	fl.helper()
	if fl.lvl >= LevelDebug {
		fl.log(append([]interface{}{"[debug]"}, v...)...)
	}
}

func (fl *funcLogger) Info(v ...interface{}) {
	fl.helper()
	if fl.lvl >= LevelInfo {
		fl.log(append([]interface{}{"[info]"}, v...)...)
	}
}

func (fl *funcLogger) Warning(v ...interface{}) {
	fl.helper()
	if fl.lvl >= LevelWarning {
		fl.log(append([]interface{}{"[warning]"}, v...)...)
	}
}

func (fl *funcLogger) Error(v ...interface{}) {
	fl.helper()
	fl.log(append([]interface{}{"[error]"}, v...)...)
}

func (fl *funcLogger) log(v ...interface{}) {
	fl.helper()

	for fn, fv := range fl.fields {
		if isEmptyFieldValue(fv) {
			continue
		}

		v = append(v, fmt.Sprintf("{%s: %v}", fn, fv))
	}

	fl.logFunc(v...)
}

func (fl *funcLogger) WithField(fn string, fv interface{}) Logger {
	fields := make(map[string]interface{})
	for cfn, cfv := range fl.fields {
		fields[cfn] = cfv
	}

	fields[fn] = fv

	return &funcLogger{
		lvl:     fl.lvl,
		logFunc: fl.logFunc,
		helper:  fl.helper,
		fields:  fields,
	}
}

func (fl *funcLogger) WithFields(newFields Fields) Logger {
	fields := make(map[string]interface{})
	for fn, fv := range fl.fields {
		fields[fn] = fv
	}
	for fn, fv := range newFields {
		fields[fn] = fv
	}

	return &funcLogger{
		lvl:     fl.lvl,
		logFunc: fl.logFunc,
		helper:  fl.helper,
		fields:  fields,
	}
}

func (fl *funcLogger) WithError(err error) Logger {
	return fl.WithField("error", err)
}

func isEmptyFieldValue(fv interface{}) bool {
	switch fvt := fv.(type) {
	case string:
		return fvt == ""

	default:
		return false
	}
}
