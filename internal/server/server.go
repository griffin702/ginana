package server

import (
	"ginana/internal/config"
	"ginana/library/conf/paladin"
	"ginana/library/log"
	"github.com/kataras/iris/v12"
	"net/http"
	"time"
)

func NewHttpServer(e *iris.Application, cfg *config.Config) (h *http.Server, err error) {
	if err = paladin.Get("http.toml").UnmarshalTOML(cfg); err != nil {
		return
	}
	if err = e.Build(); err != nil {
		log.Println(err.Error())
	}
	h = &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      e,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout),
	}
	log.Printf("HTTP服务已启动 [ http://%s ]", cfg.Server.Addr)
	err = h.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Println(err.Error())
	}
	return
}
