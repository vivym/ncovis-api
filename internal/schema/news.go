package schema

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/vivym/ncovis-api/internal/model"
)

var wordType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Word",
	Description: "Word info",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"fontSize": &graphql.Field{
			Type: graphql.Float,
		},
		"color": &graphql.Field{
			Type: graphql.String,
		},
		"rotate": &graphql.Field{
			Type: graphql.Float,
		},
		"transX": &graphql.Field{
			Type: graphql.Float,
		},
		"transY": &graphql.Field{
			Type: graphql.Float,
		},
		"fillX": &graphql.Field{
			Type: graphql.Float,
		},
		"fillY": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

var tagType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Tag",
	Description: "News Tag",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"count": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var newsType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "News",
	Description: "News Metadata",
	Fields: graphql.Fields{
		"region": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.DateTime,
		},
		"tags": &graphql.Field{
			Type: graphql.NewList(tagType),
		},
		"keywords": &graphql.Field{
			Type: graphql.NewList(wordType),
		},
		"fillingWords": &graphql.Field{
			Type: graphql.NewList(wordType),
		},
	},
})

var newsQuery = graphql.Field{
	Name:        "NewsQuery",
	Description: "News Query",
	Type: graphql.NewObject(graphql.ObjectConfig{
		Name:        "NewsQueryType",
		Description: "News Query Type",
		Fields: graphql.Fields{
			"news": &graphql.Field{
				Type: graphql.NewList(newsType),
			},
			"paging": &graphql.Field{
				Type: pagingType,
			},
		},
	}),
	Args: graphql.FieldConfigArgument{
		"region": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"date": &graphql.ArgumentConfig{
			Type: graphql.DateTime,
		},
		"cursor": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"limit": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		region, _ := p.Args["region"].(string)
		date, _ := p.Args["date"].(time.Time)
		cursor, _ := p.Args["cursor"].(string)
		limit, _ := p.Args["limit"].(int)

		if region == "" {
			return nil, errors.New("argument region is required.")
		}
		if limit == 0 || limit > 5 {
			limit = 5
		}

		from := time.Time{}
		cursorByte, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, err
		}
		_ = from.UnmarshalText(cursorByte)

		return (&model.News{}).Query(region, date, from, int64(limit))
	},
}
