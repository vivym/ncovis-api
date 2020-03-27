package schema

import (
	"encoding/base64"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/vivym/ncovis-api/internal/model"
)

var weiboTopicType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "WeiboTopic",
	Description: "Weibo Topic",
	Fields: graphql.Fields{
		"heat": &graphql.Field{
			Type: graphql.Int,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"keywords": &graphql.Field{
			Type: graphql.NewList(keywordType),
		},
	},
})

var weiboHotTopicsType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "WeiboHotTopics",
	Description: "Weibo Hot Topics",
	Fields: graphql.Fields{
		"time": &graphql.Field{
			Type: graphql.Int,
		},
		"keywords": &graphql.Field{
			Type: graphql.NewList(keywordType),
		},
		"topics": &graphql.Field{
			Type: graphql.NewList(weiboTopicType),
		},
	},
})

var weiboQuery = graphql.Field{
	Name:        "WeiboQuery",
	Description: "Weibo Query",
	Type: graphql.NewObject(graphql.ObjectConfig{
		Name:        "WeiboQueryType",
		Description: "Weibo Query Type",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(weiboHotTopicsType),
			},
			"paging": &graphql.Field{
				Type: pagingType,
			},
		},
	}),
	Args: graphql.FieldConfigArgument{
		"time": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"from": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"cursor": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"limit": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		time, _ := p.Args["time"].(int)
		from, _ := p.Args["from"].(int)
		cursor, _ := p.Args["cursor"].(string)
		limit, _ := p.Args["limit"].(int)

		/*
			if limit == 0 || limit > 10 {
				limit = 10
			}
		*/

		cursorByte, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, err
		}
		cursorInt, _ := strconv.ParseInt(string(cursorByte), 16, 64)
		if cursorInt != 0 {
			from = int(cursorInt)
		}

		return (&model.WeiboHotTopics{}).Query(int32(time), int32(from), int64(limit))
	},
}
