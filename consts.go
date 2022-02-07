package gmicro

import (
	"github.com/sirupsen/logrus"
	"github.com/yuwnloyblog/gmicro/logs"
)

func SetLogger(logger *logrus.Logger) {
	if logger != nil {
		logs.Logger = logger
	}
}
