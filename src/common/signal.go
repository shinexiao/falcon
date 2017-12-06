package common

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/gogap/logrus"
)

func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		logrus.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			reload()
		default:
			return
		}
	}
}

func reload() {
	//newConf, err := ReloadConfig()
	//if err != nil {
	//	log.Error("ReloadConfig() error(%v)", err)
	//	return
	//}
	//Conf = newConf
}
