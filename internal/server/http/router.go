package http

import (
	"ginana/internal/server/http/h_user"
	"ginana/library/mdw"
	"github.com/gin-gonic/gin"
)

func InitRouter(e *gin.Engine, u *h_user.HUser) {
	// Logger, Recovery
	e.Use(mdw.Logger, mdw.Recovery)
	// Cors
	e.Use(mdw.CORS([]string{"*"}))
	// Swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.GET("/swagger/*any", handle)

	e.GET("/", u.GetUsers)
}
