// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"ginana/internal/config"
	"ginana/internal/db"
	"ginana/internal/server"
	"ginana/internal/server/controller/api"
	"ginana/internal/server/router"
	"ginana/internal/service"
	"github.com/google/wire"
)

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewMC)
var svcProvider = wire.NewSet(service.NewErrHelper, service.New, db.NewCasbin)
var cProvider = wire.NewSet(api.New)
var httpProvider = wire.NewSet(router.InitRouter, server.NewHttpServer)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		initProvider,
		svcProvider,
		cProvider,
		httpProvider,
		NewApp,
	))
}
