package schema

import (
	"encoding/base64"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/vivym/ncovis-api/internal/model"
)

var keywordType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Keyword",
	Description: "Keyword",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"weight": &graphql.Field{
			Type: graphql.Float,
		},
		"pos": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var zhihuTopicType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ZhihuTopic",
	Description: "Zhihu Topic",
	Fields: graphql.Fields{
		"heat": &graphql.Field{
			Type: graphql.Int,
		},
		"qid": &graphql.Field{
			Type: graphql.Int,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"excerpt": &graphql.Field{
			Type: graphql.String,
		},
		"keywords": &graphql.Field{
			Type: graphql.NewList(keywordType),
		},
	},
})

var zhihuHotTopicsType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ZhihuHotTopics",
	Description: "Zhihu Hot Topics",
	Fields: graphql.Fields{
		"since": &graphql.Field{
			Type: graphql.Int,
		},
		"time": &graphql.Field{
			Type: graphql.Int,
		},
		"keywords": &graphql.Field{
			Type: graphql.NewList(keywordType),
		},
		"topics": &graphql.Field{
			Type: graphql.NewList(zhihuTopicType),
		},
	},
})

var zhihuQuery = graphql.Field{
	Name:        "ZhihuQuery",
	Description: "Zhihu Query",
	Type: graphql.NewObject(graphql.ObjectConfig{
		Name:        "ZhihuQueryType",
		Description: "Zhihu Query Type",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(zhihuHotTopicsType),
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
		"numWords": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		time, _ := p.Args["time"].(int)
		from, _ := p.Args["from"].(int)
		cursor, _ := p.Args["cursor"].(string)
		limit, _ := p.Args["limit"].(int)
		numWords, _ := p.Args["numWords"].(int)

		if numWords <= 0 || numWords > 40 {
			numWords = 40
		}

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

		return (&model.ZhihuHotTopics{}).Query(int32(time), int32(from), int64(limit), numWords)
	},
}
