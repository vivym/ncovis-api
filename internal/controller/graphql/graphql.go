package graphql

import (
  "github.com/gin-gonic/gin"
  "github.com/graphql-go/handler"

  "github.com/vivym/ncovis-api/internal/schema"
)

func GraphqlHandler() gin.HandlerFunc{
  h := handler.New(&handler.Config{
    Schema: &schema.Schema,
    Pretty: true,
    GraphiQL: true,
  })

  return func(c *gin.Context) {
    h.ServeHTTP(c.Writer, c.Request)
  }
}
