package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vivym/ncovis-api/internal/controller/graphql"
)

func New() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	r.POST("/graphql", graphql.GraphqlHandler())
	r.GET("/graphql", graphql.GraphqlHandler())

	return r
}
