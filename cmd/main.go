package main

import (
	"flag"
	_ "ginana/docs"
	"ginana/internal/config"
	"ginana/internal/wire"
	"ginana/library/conf/paladin"
	"ginana/library/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title GiNana
// @version 1.0.0
// @description 基于gin的api服务，默认端口：8000
// @host 127.0.0.1:8000
// @BasePath /api
// @license.name Apache 2.0
// @license.url
func main() {
	flag.Parse()
	closeLog := log.Init()
	log.Info("GiNana App Start")
	cfg, err := config.GetBaseConfig()
	if err != nil {
		panic(err)
	}
	if err := paladin.Init(cfg.ConfigIsLocal, cfg.ConfigPath); err != nil {
		panic(err)
	}
	_, closeFunc, err := wire.InitApp()
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			closeLog()
			log.Info("exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
