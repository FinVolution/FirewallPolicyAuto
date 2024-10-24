package router

import "github.com/kataras/iris/v12"

// Implement cross-origin resource sharing (CORS) on the server side
func cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.StatusCode(iris.StatusNoContent)
		return
	}
	ctx.Next()
}
