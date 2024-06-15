package middleware

import (
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Next()
	}
}
