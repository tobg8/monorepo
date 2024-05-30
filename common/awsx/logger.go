package awsx

import (
	awslogging "github.com/aws/smithy-go/logging"
	"github.com/monorepo/common/logging"
)

// Awslogger an abstraction for smithy-go/logging logging interface
type Awslogger struct {
	logger logging.Logger
}

// Logf logs a line on said level
func (l *Awslogger) Logf(classification awslogging.Classification, format string, v ...interface{}) {
	switch classification {
	case awslogging.Warn:
		l.logger.Warningf(format, v...)
		break
	case awslogging.Debug:
		l.logger.Debugf(format, v...)
		break
	default:
		l.logger.Errorf(format, v...)
		break
	}
}

// AwsLoggerfromLogger returns a logger compatible with aws sdk from logging.logger interface
func AwsLoggerfromLogger(logger logging.Logger) awslogging.Logger {
	return &Awslogger{logger: logger}
}
