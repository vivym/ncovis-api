package model

import (
	"encoding/base64"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Comment struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string    `json:"id" bson:",omitempty"`
	Nickname         string    `json:"nickname" bson:"nickname"`
	Title            string    `json:"title" bson:"title"`
	Desc             string    `json:"desc" bson:"desc"`
	URL              string    `json:"url" bson:"url"`
	DeviceID         string    `json:"deviceId" bson:"deviceId"`
	IsPublished      bool      `json:"isPublished" bson:"isPublished"`
	IsTop            bool      `json:"isTop" bson:"isTop"`
	IsDeleted        bool      `json:"isDeleted" bson:"isDeleted"`
	IsMine           bool      `json:"isMine" bson:",omitempty"`
	ViewCount        int32     `json:"viewCount" bson:"viewCount"`
	Keywords         []Keyword `json:"keywords" bson:"keywords"`
}

type commentQueryResult struct {
	Data   []*Comment `json:"data"`
	Paging Paging     `json:"paging"`
}

func (c *Comment) prepareExtraFields(deviceID string) {
	id, _ := c.GetID().(primitive.ObjectID)
	c.ID = id.Hex()
	c.IsMine = deviceID == c.DeviceID
}

type Success struct {
	Code int `json:"code"`
}

func (*Comment) Query(sortBy string, offset, limit int, isAdmin bool, deviceID string) (*commentQueryResult, error) {
	extraComments := []*Comment{}
	if offset == 0 {
		filter := bson.M{
			"isDeleted":   false,
			"isPublished": false,
		}
		findOptions := options.Find()
		findOptions.SetSort(bson.M{"created_at": -1})
		if !isAdmin {
			filter["deviceId"] = deviceID
		}
		if err := mgm.Coll(&Comment{}).SimpleFind(&extraComments, filter, findOptions); err != nil {
			return nil, err
		}
	}

	filter := bson.M{}
	filter["isDeleted"] = false
	filter["isPublished"] = true

	var err error
	var count int64
	count, err = mgm.Coll(&Comment{}).CountDocuments(mgm.Ctx(), filter)
	if err != nil {
		return nil, err
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"isTop": -1, sortBy: 1})

	comments := []*Comment{}
	if err = mgm.Coll(&Comment{}).SimpleFind(&comments, filter, findOptions); err != nil {
		return nil, err
	}
	comments = append(extraComments, comments...)
	for _, comment := range comments {
		comment.prepareExtraFields(deviceID)
	}

	nextCursor := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(offset + len(comments))))

	result := commentQueryResult{
		Data: comments,
		Paging: Paging{
			Total:      count + int64(len(extraComments)),
			NextCursor: nextCursor,
		},
	}
	return &result, err
}

func (*Comment) QueryWithID(id string) (*Comment, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": _id,
	}

	comments := []*Comment{}
	if err := mgm.Coll(&Comment{}).SimpleFind(&comments, filter); err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errors.New("not found")
	}
	return comments[0], nil
}

func (*Comment) Create(nickname, title, desc, url, deviceID string) (*Comment, error) {
	comment := Comment{
		Nickname:    nickname,
		Title:       title,
		Desc:        desc,
		URL:         url,
		DeviceID:    deviceID,
		IsPublished: false,
		IsTop:       false,
		IsDeleted:   false,
		ViewCount:   0,
	}

	if err := mgm.Coll(&Comment{}).Create(&comment); err != nil {
		return nil, err
	}
	comment.prepareExtraFields(deviceID)

	return &comment, nil
}

func (*Comment) Publish(id string, isPublish bool) (*Success, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$set": bson.M{
			"isPublished": isPublish,
		},
	}
	_, err = mgm.Coll(&Comment{}).UpdateOne(mgm.Ctx(), filter, update)
	if err != nil {
		return nil, err
	}

	return &Success{}, nil
}

func (*Comment) Top(id string, isTop bool) (*Success, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$set": bson.M{
			"isTop": isTop,
		},
	}
	_, err = mgm.Coll(&Comment{}).UpdateOne(mgm.Ctx(), filter, update)
	if err != nil {
		return nil, err
	}

	return &Success{}, nil
}

func (*Comment) View(id string) (*Success, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$inc": bson.M{
			"viewCount": 1,
		},
	}
	_, err = mgm.Coll(&Comment{}).UpdateOne(mgm.Ctx(), filter, update)
	if err != nil {
		return nil, err
	}

	return &Success{}, nil
}

func (*Comment) Delete(id string) (*Success, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}
	_, err = mgm.Coll(&Comment{}).UpdateOne(mgm.Ctx(), filter, update)
	if err != nil {
		return nil, err
	}

	return &Success{}, nil
}
