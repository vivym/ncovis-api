package model

import (
	"encoding/base64"
	"strconv"

	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Keyword struct {
	Name   string  `json:"name" bson:"name"`
	Weight float64 `json:"weight" bson:"weight"`
	POS    string  `json:"pos" bson:"pos"`
}

type ZhihuHotTopics struct {
	mgm.DefaultModel `bson:",inline"`
	Since            int32     `json:"since" bson:"since"`
	Time             int32     `json:"time" bson:"time"`
	Keywords         []Keyword `json:"keywords" bson:"keywords"`
	Topics           []struct {
		Heat     int32     `json:"heat" bson:"heat"`
		QID      int32     `json:"qid" bson:"qid"`
		Title    string    `json:"title" bson:"title"`
		Excerpt  string    `json:"excerpt" bson:"excerpt"`
		Keywords []Keyword `json:"keywords" bson:"keywords"`
	} `json:"topics" bson:"topics"`
}

func (z *ZhihuHotTopics) CollectionName() string {
	return "zhihu_hot_topics"
}

type zhihuQueryResult struct {
	Data   []*ZhihuHotTopics `json:"data"`
	Paging Paging            `json:"paging"`
}

func (*ZhihuHotTopics) Query(time, from int32, limit int64) (*zhihuQueryResult, error) {
	filter := bson.M{}
	if time != 0 {
		filter["time"] = time
	}

	var err error
	var count int64
	count, err = mgm.Coll(&ZhihuHotTopics{}).CountDocuments(mgm.Ctx(), filter)
	if err != nil {
		return nil, err
	}

	if from != 0 {
		filter["time"] = bson.M{"$gt": from}
	}

	findOptions := options.Find()
	if limit != 0 {
		findOptions.SetLimit(limit)
	}
	findOptions.SetSort(bson.M{"time": 1})

	data := []*ZhihuHotTopics{}
	if err = mgm.Coll(&ZhihuHotTopics{}).SimpleFind(&data, filter, findOptions); err != nil {
		return nil, err
	}

	lastDate := ""
	if len(data) > 0 {
		lastDate = strconv.FormatInt(int64(data[len(data)-1].Time), 16)
	}
	nextCursor := base64.StdEncoding.EncodeToString([]byte(lastDate))

	for _, topics := range data {
		count := len(topics.Keywords)
		if count > 40 {
			count = 40
		}
		topics.Keywords = topics.Keywords[:count]
		for _, topic := range topics.Topics {
			count := len(topic.Keywords)
			if count > 40 {
				count = 40
			}
			topic.Keywords = topic.Keywords[:count]
		}
	}

	result := zhihuQueryResult{
		Data: data,
		Paging: Paging{
			Total:      count,
			NextCursor: nextCursor,
		},
	}
	return &result, err
}
