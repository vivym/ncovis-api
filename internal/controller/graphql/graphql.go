package graphql

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"

	"github.com/vivym/ncovis-api/internal/schema"
)

func GraphqlHandler(adminToken string) gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:   &schema.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		deviceID := c.Request.Header.Get("DeviceID")
		if len(token) > 7 {
			token = token[7:]
		}
		ctx := context.WithValue(c.Request.Context(), "user", struct {
			IsAdmin  bool
			DeviceID string
		}{
			IsAdmin:  token == adminToken,
			DeviceID: deviceID,
		})
		h.ContextHandler(ctx, c.Writer, c.Request)
	}
}
