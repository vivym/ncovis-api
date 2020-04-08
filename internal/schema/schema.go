package schema

import "github.com/graphql-go/graphql"

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name:        "RootQuery",
	Description: "Root Query",
	Fields: graphql.Fields{
		"ncov":        &ncovQuery,
		"news":        &newsQuery,
		"zhihu":       &zhihuQuery,
		"weibo":       &weiboQuery,
		"ncovOverall": &ncovOverallQuery,
		"comment":     &commentQuery,
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name:        "RootMutation",
	Description: "Root Mutation",
	Fields: graphql.Fields{
		"createComment":  &createComment,
		"publishComment": &publishComment,
		"topComment":     &topComment,
		"viewComment":    &viewComment,
		"deleteComment":  &deleteComment,
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
