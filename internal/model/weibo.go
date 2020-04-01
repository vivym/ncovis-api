package model

import (
	"encoding/base64"
	"strconv"

	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WeiboHotTopics struct {
	mgm.DefaultModel `bson:",inline"`
	Time             int32        `json:"time" bson:"time"`
	Keywords         []Keyword    `json:"keywords" bson:"keywords"`
	Topics           []WeiboTopic `json:"topics" bson:"topics"`
}

type WeiboTopic struct {
	Heat     int32     `json:"heat" bson:"heat"`
	URL      string    `json:"url" bson:"url"`
	Title    string    `json:"title" bson:"title"`
	Keywords []Keyword `json:"keywords" bson:"keywords"`
}

func (z *WeiboHotTopics) CollectionName() string {
	return "weibo_hot_topics"
}

type weiboQueryResult struct {
	Data   []*WeiboHotTopics `json:"data"`
	Paging Paging            `json:"paging"`
}

func (*WeiboHotTopics) Query(time, from int32, limit int64, numWords int) (*weiboQueryResult, error) {
	filter := bson.M{}
	if time != 0 {
		filter["time"] = time
	}

	var err error
	var count int64
	count, err = mgm.Coll(&WeiboHotTopics{}).CountDocuments(mgm.Ctx(), filter)
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

	data := []*WeiboHotTopics{}
	if err = mgm.Coll(&WeiboHotTopics{}).SimpleFind(&data, filter, findOptions); err != nil {
		return nil, err
	}

	lastDate := ""
	if len(data) > 0 {
		lastDate = strconv.FormatInt(int64(data[len(data)-1].Time), 16)
	}
	nextCursor := base64.StdEncoding.EncodeToString([]byte(lastDate))

	for _, topics := range data {
		count := len(topics.Keywords)
		if count > numWords {
			count = numWords
		}
		topics.Keywords = topics.Keywords[:count]
		for _, topic := range topics.Topics {
			count := len(topic.Keywords)
			if count > numWords {
				count = numWords
			}
			topic.Keywords = topic.Keywords[:count]
		}
	}

	result := weiboQueryResult{
		Data: data,
		Paging: Paging{
			Total:      count,
			NextCursor: nextCursor,
		},
	}
	return &result, err
}
