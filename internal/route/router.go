package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vivym/ncovis-api/internal/controller/graphql"
)

func New(graphiQLToken string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	r.POST("/graphql", graphql.GraphqlHandler())
	r.GET("/graphql", func(c *gin.Context) {
		token := c.Query("token")
		if token != graphiQLToken {
			c.String(403, "Permission Diend.")
			c.Abort()
			return
		}
		c.Next()
	}, graphql.GraphqlHandler())

	return r
}
