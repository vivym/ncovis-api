package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vivym/ncovis-api/internal/controller/graphql"
)

func New() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	r.POST("graphql", graphql.GraphqlHandler())
	r.GET("graphql", graphql.GraphqlHandler())

	return r
}
