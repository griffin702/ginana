package api

import (
	"ginana/internal/model"
	"ginana/internal/server/resp"
	"ginana/internal/service"
	"github.com/kataras/iris/v12"
)

type CApi struct {
	//Ctx iris.Context
	svc *service.Service
}

func New(s *service.Service) *CApi {
	return &CApi{
		svc: s,
	}
}

// GetUsers godoc
// @Description 获取用户列表(分页)
// @Tags Users
// @Accept  json
// @Produce  json
// @Param page query int true "页码"
// @Param pagesize query int true "页码尺寸"
// @Success 200 {object} model.User
// @Failure 500 {object} resp.JSON
// @Router /users [get]
func (c *CApi) GetUsers(ctx iris.Context) {
	data := model.GiNana{
		Hello: "Hello GiNana!",
	}
	ctx.JSON(resp.PlusJson(data, nil))
}
