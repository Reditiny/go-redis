package mylog

import (
	"github.com/sirupsen/logrus"
	"redis-based-on-go/config"
)

var Logger = logrus.New()

func InitLog() {
	Logger.Level = logrus.Level(config.Conf.Logger.Level)
}
