package abortwithjson

import "github.com/gin-gonic/gin"

func main2() {
	router := gin.Default()
	cc := gin.Context{}
	router.GET("/", Get)
	router.Run()
	cc.AbortWithStatusJSON(200, gin.H{"message": "Hello world"})
}

func Set(ginCtx *gin.Context) {
	ginCtx.AbortWithStatusJSON(200, gin.H{"message": "Hello world"})
}

func Get(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(200, gin.H{"message": "Hello world"})
}
