package logs

import "github.com/sirupsen/logrus"

func Panic(f interface{}, v ...interface{}) {
	logrus.Panic(f, v)
}

func Fata(f interface{}, v ...interface{}) {
	logrus.Fatal(f, v)
}

func Error(f interface{}, v ...interface{}) {
	logrus.Error(f, v)
}
func Warn(f interface{}, v ...interface{}) {
	logrus.Warn(f, v)
}

func Info(f interface{}, v ...interface{}) {
	logrus.Info(f, v)
}

func Debug(f interface{}, v ...interface{}) {
	logrus.Debug(f, v)
}

func Trace(f interface{}, v ...interface{}) {
	logrus.Trace(f, v)
}
