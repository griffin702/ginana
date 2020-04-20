package router

import (
	"ginana/internal/config"
	"ginana/internal/server/controller/api"
	"ginana/library/mdw"
	"github.com/kataras/iris/v12"
)

func InitRouter(api *api.CApi, cfg *config.Config) (e *iris.Application, err error) {
	e = NewIris(cfg)
	//sessManager := sessions.New(sessions.Config{
	//	Cookie:  "GiNana_Session",
	//	Expires: 24 * time.Hour,
	//})
	apiParty := e.Party("/api", mdw.CORS([]string{"*"})).AllowMethods(iris.MethodOptions)
	{
		apiParty.Get("/", api.GetUsers)
	}
	return
}
