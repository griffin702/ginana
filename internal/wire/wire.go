// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"ginana/internal/config"
	"ginana/internal/db"
	"ginana/internal/server/http"
	"ginana/internal/service"
	"github.com/google/wire"
)

var serviceProvider = wire.NewSet(config.NewConfig, db.NewDB, service.New)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		serviceProvider,
		http.NewGin,
		http.NewHttpServer,
		db.NewCasbin,
		NewApp,
	))
}
