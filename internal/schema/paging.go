package schema

import "github.com/graphql-go/graphql"

var pagingType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Paging",
	Description: "Paging Info",
	Fields: graphql.Fields{
		"total": &graphql.Field{
			Type: graphql.Int,
		},
		"nextCursor": &graphql.Field{
			Type: graphql.String,
		},
	},
})
