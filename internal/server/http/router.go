package http

import (
	"ginana/internal/server/http/h_user"
	"github.com/gin-gonic/gin"
)

func InitRouter(u *h_user.HUser) (e *gin.Engine) {
	e = NewGin()

	e.GET("/", u.GetUsers)

	return
}
