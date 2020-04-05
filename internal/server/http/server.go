package http

import (
	"ginana/internal/config"
	"ginana/library/conf/paladin"
	"ginana/library/log"
	"ginana/library/mdw"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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

func NewGin() (e *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.GetOutFile()
	e = gin.Default()
	// Logger, Recovery
	e.Use(mdw.Logger, mdw.Recovery)
	// Cors
	e.Use(mdw.CORS([]string{"*"}))
	// Swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.GET("/swagger/*any", handle)
	return
}
