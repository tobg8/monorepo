package logrus

import (
	"github.com/sirupsen/logrus"
)

type appNameHook struct {
	appName string
}

func (*appNameHook) Levels() []logrus.Level { return logrus.AllLevels }

func (h *appNameHook) Fire(entry *logrus.Entry) error {
	if _, ok := entry.Data["application"]; !ok {
		entry.Data["application"] = h.appName
	}

	return nil
}
