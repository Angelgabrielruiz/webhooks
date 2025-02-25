package infrastructure

import (
	"github_wb/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	routes := engine.Group("pull_request")

	{
		routes.POST("process", handlers.PullRequestEvent)
	}
 //comentariogggggwww
}
