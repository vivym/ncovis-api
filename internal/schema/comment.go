package schema

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/vivym/ncovis-api/internal/nlp"

	"github.com/graphql-go/graphql"
	"github.com/vivym/ncovis-api/internal/model"
)

var commentType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Comment",
	Description: "Comment",
	Fields: graphql.Fields{
		"code": &graphql.Field{
			Type: graphql.Int,
		},
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"nickname": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"desc": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"isPublished": &graphql.Field{
			Type: graphql.Boolean,
		},
		"isTop": &graphql.Field{
			Type: graphql.Boolean,
		},
		"isMine": &graphql.Field{
			Type: graphql.Boolean,
		},
		"viewCount": &graphql.Field{
			Type: graphql.Int,
		},
		"keywords": &graphql.Field{
			Type: graphql.NewList(keywordType),
		},
	},
})

var commentQuery = graphql.Field{
	Name:        "CommentQuery",
	Description: "Comment Query",
	Type: graphql.NewObject(graphql.ObjectConfig{
		Name:        "CommentQueryType",
		Description: "Comment Query Type",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(commentType),
			},
			"paging": &graphql.Field{
				Type: pagingType,
			},
		},
	}),
	Args: graphql.FieldConfigArgument{
		"sortBy": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"cursor": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"limit": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		sortBy, _ := p.Args["sortBy"].(string)
		cursor, _ := p.Args["cursor"].(string)
		limit, _ := p.Args["limit"].(int)
		user, _ := p.Context.Value("user").(struct {
			IsAdmin  bool
			DeviceID string
		})

		if sortBy == "" || (sortBy != "viewCount" && sortBy != "created_at") {
			sortBy = "viewCount"
		}
		if limit == 0 || limit > 20 {
			limit = 20
		}

		cursorByte, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, err
		}
		offset, _ := strconv.Atoi(string(cursorByte))

		return (&model.Comment{}).Query(sortBy, offset, limit, user.IsAdmin, user.DeviceID)
	},
}

var createComment = graphql.Field{
	Name:        "CreateComment",
	Description: "Create Comment",
	Type:        commentType,
	Args: graphql.FieldConfigArgument{
		"nickname": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"desc": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"url": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"numWords": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		nickname, _ := p.Args["nickname"].(string)
		title, _ := p.Args["title"].(string)
		desc, _ := p.Args["desc"].(string)
		url, _ := p.Args["url"].(string)
		numWords, _ := p.Args["numWords"].(int)

		if numWords <= 0 || numWords > 40 {
			numWords = 40
		}

		user, _ := p.Context.Value("user").(struct {
			IsAdmin  bool
			DeviceID string
		})

		nickname = strings.TrimSpace(nickname)
		title = strings.TrimSpace(title)
		desc = strings.TrimSpace(desc)
		url = strings.TrimSpace(url)

		if len(nickname) == 0 || len(nickname) > 256 {
			return nil, errors.New("invalid nickname")
		}

		if len(title) == 0 || len(title) > 256 {
			return nil, errors.New("invalid title")
		}

		if len(desc) > 1024*1024 {
			return nil, errors.New("invalid desc")
		}

		if len(url) > 1024 {
			return nil, errors.New("invalid url")
		}

		if len(url) > 0 && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return nil, errors.New("invalid url")
		}

		keywords, _ := nlp.Get().ExtractKeywords(title+"\n"+desc, int64(numWords))

		return (&model.Comment{}).Create(nickname, title, desc, url, user.DeviceID, keywords)
	},
}

var publishComment = graphql.Field{
	Name:        "PublishComment",
	Description: "Publish Comment",
	Type:        commentType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"isPublish": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, _ := p.Args["id"].(string)
		isPublish, _ := p.Args["isPublish"].(bool)
		user, _ := p.Context.Value("user").(struct {
			IsAdmin  bool
			DeviceID string
		})

		if !user.IsAdmin {
			comment, err := (&model.Comment{}).QueryWithID(id)
			if err != nil || comment.DeviceID != user.DeviceID {
				return nil, errors.New("Permission Diend.")
			}
		}

		return (&model.Comment{}).Publish(id, isPublish)
	},
}

var topComment = graphql.Field{
	Name:        "TopComment",
	Description: "Top Comment",
	Type:        commentType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"isTop": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, _ := p.Args["id"].(string)
		isTop, _ := p.Args["isTop"].(bool)
		user, _ := p.Context.Value("user").(struct {
			IsAdmin  bool
			DeviceID string
		})

		if !user.IsAdmin {
			return nil, errors.New("Permission Diend.")
		}

		return (&model.Comment{}).Top(id, isTop)
	},
}

var viewComment = graphql.Field{
	Name:        "ViewComment",
	Description: "View Comment",
	Type:        commentType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, _ := p.Args["id"].(string)

		return (&model.Comment{}).View(id)
	},
}

var deleteComment = graphql.Field{
	Name:        "DeleteComment",
	Description: "Delete Comment",
	Type:        commentType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, _ := p.Args["id"].(string)
		user, _ := p.Context.Value("user").(struct {
			IsAdmin  bool
			DeviceID string
		})

		if !user.IsAdmin {
			comment, err := (&model.Comment{}).QueryWithID(id)
			if err != nil || comment.DeviceID != user.DeviceID {
				return nil, errors.New("Permission Diend.")
			}
		}

		return (&model.Comment{}).Delete(id)
	},
}
