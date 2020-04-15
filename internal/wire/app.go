package wire

import (
	"context"
	"ginana/internal/service"
	"ginana/library/log"
	"github.com/casbin/casbin/v2"
	"net/http"
	"time"
)

type App struct {
	Server *http.Server
	svc    service.Service
	ef     *casbin.SyncedEnforcer
}

func NewApp(server *http.Server, svc service.Service, ef *casbin.SyncedEnforcer) (app *App, closeFunc func(), err error) {
	app = &App{
		Server: server,
		svc:    svc,
		ef:     ef,
	}
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			log.Errorf("httpServer.Shutdown error(%v)", err)
		}
		svc.Close()
		cancel()
	}
	return
}
