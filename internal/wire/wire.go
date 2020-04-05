// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"ginana/internal/config"
	"ginana/internal/db"
	"ginana/internal/server/http"
	"ginana/internal/server/http/h_user"
	"ginana/internal/service"
	"ginana/internal/service/user"
	"github.com/google/wire"
)

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewCasbin)
var iProvider = wire.NewSet(user.New)
var hProvider = wire.NewSet(h_user.New)
var httpProvider = wire.NewSet(http.InitRouter, http.NewHttpServer)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		initProvider,
		iProvider,
		service.New,
		hProvider,
		httpProvider,
		NewApp,
	))
}
