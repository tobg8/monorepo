package logrus

import (
	"github.com/sirupsen/logrus"
)

type versionHook struct {
	version string
}

func (*versionHook) Levels() []logrus.Level { return logrus.AllLevels }

func (h *versionHook) Fire(entry *logrus.Entry) error {
	entry.Data["@version"] = h.version
	return nil
}
