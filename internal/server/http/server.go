package http

import (
	"ginana/internal/config"
	"ginana/internal/server/http/h_user"
	"ginana/library/conf/paladin"
	"ginana/library/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NewGin(u *h_user.HUser) (e *gin.Engine, err error) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.GetOutFile()
	e = gin.Default()
	InitRouter(e, u)
	return
}

func NewHttpServer(e *gin.Engine, cfg *config.Config) (h *http.Server, err error) {
	if err = paladin.Get("http.toml").UnmarshalTOML(cfg); err != nil {
		return
	}
	h = &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      e,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout),
	}
	log.Infof("HTTP服务已启动 [ http://%s ]", cfg.Server.Addr)
	err = h.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorf(err.Error())
	}
	return
}
