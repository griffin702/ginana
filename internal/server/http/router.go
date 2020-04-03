package http

import (
	"ginana/library/mdw"
	"github.com/gin-gonic/gin"
)

func InitRouter(e *gin.Engine) {
	// Logger, Recovery
	e.Use(mdw.Logger, mdw.Recovery)
	// Cors
	e.Use(mdw.CORS([]string{"*"}))
	// Swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.GET("/swagger/*any", handle)

	e.GET("/", howToStart)
}
