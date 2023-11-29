package router

import (
	_tagHandler "go-postgres-clean-arch/tag/delivery/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(tagHandler *_tagHandler.TagHandler) *gin.Engine {
	router := gin.Default()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome home")
	})

	baseRouter := router.Group("/api")
	tagsRouter := baseRouter.Group("/tags")
	tagsRouter.GET("", tagHandler.FindByAll)
	tagsRouter.GET("/:tagId", tagHandler.FindById)
	tagsRouter.POST("", tagHandler.Create)
	tagsRouter.PATCH("/:tagId", tagHandler.Update)
	tagsRouter.DELETE("/:tagId", tagHandler.Delete)

	return router
}
