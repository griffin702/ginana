package hook

import (
	"github.com/sirupsen/logrus"
)

type DefaultFieldHook struct {
}

func (h *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["appName"] = "GiNana"
	return nil
}

func (h *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
