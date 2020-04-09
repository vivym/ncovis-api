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
		"country": &graphql.Field{
			Type: graphql.String,
		},
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

var ncovOverallType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "nCoVGlobalInfo",
	Description: "nCoV Global info",
	Fields: graphql.Fields{
		"dead": &graphql.Field{
			Type: graphql.Int,
		},
		"deadIncr": &graphql.Field{
			Type: graphql.Int,
		},
		"confirmed": &graphql.Field{
			Type: graphql.Int,
		},
		"confirmedIncr": &graphql.Field{
			Type: graphql.Int,
		},
		"cured": &graphql.Field{
			Type: graphql.Int,
		},
		"curedIncr": &graphql.Field{
			Type: graphql.Int,
		},
		"remainingConfirmed": &graphql.Field{
			Type: graphql.Int,
		},
		"remainingConfirmedIncr": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var ncovQuery = graphql.Field{
	Name:        "NCoVQuery",
	Description: "nCoV Info Query",
	Type:        graphql.NewList(ncovType),
	Args: graphql.FieldConfigArgument{
		"country": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"region": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"date": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		country, _ := p.Args["country"].(string)
		region, _ := p.Args["region"].(string)
		date, _ := p.Args["date"].(string)

		return (&model.NCoVInfo{}).Query(country, region, date)
	},
}

var ncovOverallQuery = graphql.Field{
	Name:        "NCoVOverallQuery",
	Description: "nCoV Overall Info Query",
	Type: graphql.NewObject(graphql.ObjectConfig{
		Name:        "nCoVOverall",
		Description: "nCoV overall info",
		Fields: graphql.Fields{
			"time": &graphql.Field{
				Type: graphql.Int,
			},
			"dead": &graphql.Field{
				Type: graphql.Int,
			},
			"deadIncr": &graphql.Field{
				Type: graphql.Int,
			},
			"confirmed": &graphql.Field{
				Type: graphql.Int,
			},
			"confirmedIncr": &graphql.Field{
				Type: graphql.Int,
			},
			"suspected": &graphql.Field{
				Type: graphql.Int,
			},
			"suspectedIncr": &graphql.Field{
				Type: graphql.Int,
			},
			"cured": &graphql.Field{
				Type: graphql.Int,
			},
			"curedIncr": &graphql.Field{
				Type: graphql.Int,
			},
			"remainingConfirmed": &graphql.Field{
				Type: graphql.Int,
			},
			"remainingConfirmedIncr": &graphql.Field{
				Type: graphql.Int,
			},
			"serious": &graphql.Field{
				Type: graphql.Int,
			},
			"seriousIncr": &graphql.Field{
				Type: graphql.Int,
			},
			"global": &graphql.Field{
				Type: ncovOverallType,
			},
		},
	}),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return (&model.NCovOverallInfo{}).Query()
	},
}
