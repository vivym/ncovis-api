package model

import (
	"encoding/base64"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type News struct {
	mgm.DefaultModel `bson:",inline"`
	Region           string    `json:"region" bson:"region"`
	Date             time.Time `json:"date" bson:"date"`
	Tags             []struct {
		Name  string `json:"name" bson:"name"`
		Count int    `json:"count" bson:"count"`
	} `json:"tags" bson:"tags"`
	Keywords     []word `json:"keywords" bson:"keywords"`
	FillingWords []word `json:"fillingWords" bson:"fillingWords"`
}

type word struct {
	Name     string  `json:"name" bson:"name"`
	FontSize float64 `json:"fontSize" bson:"fontSize"`
	Color    string  `json:"color" bson:"color"`
	Rotate   float64 `json:"rotate" bson:"rotate"`
	TransX   float64 `json:"transX" bson:"transX"`
	TransY   float64 `json:"transY" bson:"transY"`
	FillX    float64 `json:"fillX" bson:"fillX"`
	FillY    float64 `json:"fillY" bson:"fillY"`
}

type Paging struct {
	Total      int64  `json:"total"`
	NextCursor string `json:"nextCursor"`
}

type newsQueryResult struct {
	News   []*News `json:"news"`
	Paging Paging  `json:"paging"`
}

func (*News) Query(region string, date, from time.Time, limit int64) (*newsQueryResult, error) {
	filter := bson.M{}
	if region != "" {
		filter["region"] = region
	}
	if !date.IsZero() {
		filter["date"] = date
	}

	var err error
	var count int64
	count, err = mgm.Coll(&News{}).CountDocuments(mgm.Ctx(), filter)
	if err != nil {
		return nil, err
	}

	if !from.IsZero() {
		filter["date"] = bson.M{"$gt": from}
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"date": 1})

	news := []*News{}
	if err = mgm.Coll(&News{}).SimpleFind(&news, filter, findOptions); err != nil {
		return nil, err
	}

	lastDate := []byte{}
	if len(news) > 0 {
		lastDate, _ = news[len(news)-1].Date.MarshalText()
	}
	nextCursor := base64.StdEncoding.EncodeToString(lastDate)

	result := newsQueryResult{
		News: news,
		Paging: Paging{
			Total:      count,
			NextCursor: nextCursor,
		},
	}
	return &result, err
}
