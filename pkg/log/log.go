package log

import (
	"github.com/aquaproj/aqua/pkg/runtime"
	"github.com/sirupsen/logrus"
)

func New(rt *runtime.Runtime, version string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"aqua_version": version,
		"program":      "aqua",
		"env":          rt.GOOS + "/" + rt.GOARCH,
	})
}

func SetLevel(level string, logE *logrus.Entry) {
	if level == "" {
		return
	}
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logE.WithField("log_level", level).WithError(err).Error("the log level is invalid")
		return
	}
	logrus.SetLevel(lvl)
}
