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
	"ginana/internal/service/i_user"
	"github.com/google/wire"
)

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewCasbin)
var iProvider = wire.NewSet(i_user.New)
var cProvider = wire.NewSet(api.New)
var httpProvider = wire.NewSet(router.InitRouter, server.NewHttpServer)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		initProvider,
		iProvider,
		service.New,
		cProvider,
		httpProvider,
		NewApp,
	))
}
