package router

import (
	"ginana/internal/config"
	"ginana/internal/server/controller/api"
	"github.com/kataras/iris/v12"
)

func InitRouter(api *api.CApi, cfg *config.Config) (e *iris.Application) {
	e = NewIris(cfg)
	//sessManager := sessions.New(sessions.Config{
	//	Cookie:  "GiNana_Session",
	//	Expires: 24 * time.Hour,
	//})
	apiParty := e.Party("/api", Cors()).AllowMethods(iris.MethodOptions)
	{
		apiParty.Get("/", api.GetUsers)
	}
	return
}
