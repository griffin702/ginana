// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"ginana/internal/config"
	"ginana/internal/db"
	"ginana/internal/server/http"
	"ginana/internal/service"
	"ginana/internal/service/user"
	"github.com/google/wire"
)

var initProvider = wire.NewSet(config.NewConfig, db.NewDB)
var serviceProvider = wire.NewSet(user.New)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		initProvider,
		serviceProvider,
		service.New,
		http.NewGin,
		http.NewHttpServer,
		db.NewCasbin,
		NewApp,
	))
}
