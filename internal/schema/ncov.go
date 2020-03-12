package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/vivym/ncovis-api/internal/model"
)

var cityType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "nCoVCity",
	Description: "nCoV city info",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"locID": &graphql.Field{
			Type: graphql.Int,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"dead": &graphql.Field{
			Type: graphql.Int,
		},
		"confirmed": &graphql.Field{
			Type: graphql.Int,
		},
		"suspected": &graphql.Field{
			Type: graphql.Int,
		},
		"cured": &graphql.Field{
			Type: graphql.Int,
		},
		"remainingConfirmed": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var ncovType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "nCoVRegion",
	Description: "nCoV region info",
	Fields: graphql.Fields{
		"region": &graphql.Field{
			Type: graphql.String,
		},
		"locID": &graphql.Field{
			Type: graphql.Int,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"dead": &graphql.Field{
			Type: graphql.Int,
		},
		"confirmed": &graphql.Field{
			Type: graphql.Int,
		},
		"suspected": &graphql.Field{
			Type: graphql.Int,
		},
		"cured": &graphql.Field{
			Type: graphql.Int,
		},
		"remainingConfirmed": &graphql.Field{
			Type: graphql.Int,
		},
		"cities": &graphql.Field{
			Type: graphql.NewList(cityType),
		},
	},
})

var ncovQuery = graphql.Field{
	Name:        "NCoVQuery",
	Description: "nCoV Info Query",
	Type:        graphql.NewList(ncovType),
	Args: graphql.FieldConfigArgument{
		"region": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"date": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		region, _ := p.Args["region"].(string)
		date, _ := p.Args["date"].(string)
		return (&model.NCoVInfo{}).Query(region, date)
	},
}
