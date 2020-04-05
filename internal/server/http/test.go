package http

import (
	"github.com/gin-gonic/gin"
)

// BlogGin hello BlogGin.
type GiNana struct {
	Hello string
}

// howToStart godoc
// @Description 测试
// @Tags Public
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} string "ok"
// @Failure 500 {string} string "failed"
// @Router / [get]
func howToStart(ctx *gin.Context) {
	k := &GiNana{
		Hello: "GiNana Server",
	}
	ctx.JSON(200, k)
}
