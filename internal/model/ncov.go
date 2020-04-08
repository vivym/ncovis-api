package model

import (
	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NCoVInfo struct {
	mgm.DefaultModel   `bson:",inline"`
	Country            string `json:"country" bson:"country"`
	Region             string `json:"region" bson:"region"`
	LocID              int    `json:"locID" bson:"locID"`
	Date               string `json:"date" bson:"date"`
	Dead               int    `json:"dead" bson:"dead"`
	Confirmed          int    `json:"confirmed" bson:"confirmed"`
	Suspected          int    `json:"suspected" bson:"suspected"`
	Cured              int    `json:"cured" bson:"cured"`
	RemainingConfirmed int    `json:"remainingConfirmed" bson:"remainingConfirmed"`
	Cities             []struct {
		Name               string `json:"name" bson:"name"`
		LocID              int    `json:"locID" bson:"locID"`
		Dead               int    `json:"dead" bson:"dead"`
		Confirmed          int    `json:"confirmed" bson:"confirmed"`
		Suspected          int    `json:"suspected" bson:"suspected"`
		Cured              int    `json:"cured" bson:"cured"`
		RemainingConfirmed int    `json:"remainingConfirmed" bson:"remainingConfirmed"`
	} `json:"cities" bson:"cities"`
}

type NCovOverallInfo struct {
	mgm.DefaultModel       `bson:",inline"`
	Time                   int32 `json:"time" bson:"time"`
	Dead                   int   `json:"dead" bson:"dead"`
	DeadIncr               int   `json:"deadIncr" bson:"deadIncr"`
	Confirmed              int   `json:"confirmed" bson:"confirmed"`
	ConfirmedIncr          int   `json:"confirmedIncr" bson:"confirmedIncr"`
	Suspected              int   `json:"suspected" bson:"suspected"`
	SuspectedIncr          int   `json:"suspectedIncr" bson:"suspectedIncr"`
	Cured                  int   `json:"cured" bson:"cured"`
	CuredIncr              int   `json:"curedIncr" bson:"curedIncr"`
	RemainingConfirmed     int   `json:"remainingConfirmed" bson:"remainingConfirmed"`
	RemainingConfirmedIncr int   `json:"remainingConfirmedIncr" bson:"remainingConfirmedIncr"`
	Serious                int   `json:"serious" bson:"serious"`
	SeriousIncr            int   `json:"seriousIncr" bson:"seriousIncr"`
	Global                 struct {
		Dead                   int `json:"dead" bson:"dead"`
		DeadIncr               int `json:"deadIncr" bson:"deadIncr"`
		Confirmed              int `json:"confirmed" bson:"confirmed"`
		ConfirmedIncr          int `json:"confirmedIncr" bson:"confirmedIncr"`
		Cured                  int `json:"cured" bson:"cured"`
		CuredIncr              int `json:"curedIncr" bson:"curedIncr"`
		RemainingConfirmed     int `json:"remainingConfirmed" bson:"remainingConfirmed"`
		RemainingConfirmedIncr int `json:"remainingConfirmedIncr" bson:"remainingConfirmedIncr"`
	} `json:"global" bson:"global"`
}

func (*NCoVInfo) CollectionName() string {
	return "ncovinfos"
}

func (n *NCovOverallInfo) CollectionName() string {
	return "ncov_overall_infos"
}

func (*NCoVInfo) Query(country, region, date string) ([]*NCoVInfo, error) {
	filter := bson.M{}
	if country != "" {
		filter["country"] = country
	}
	if region != "" {
		filter["region"] = region
	}
	if date != "" {
		filter["date"] = date
	}
	infos := []*NCoVInfo{}
	if err := mgm.Coll(&NCoVInfo{}).SimpleFind(&infos, filter); err != nil {
		return infos, err
	}
	return infos, nil
}

func (*NCovOverallInfo) Query() (*NCovOverallInfo, error) {
	infos := []*NCovOverallInfo{}

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetLimit(1)
	findOptions.SetSort(bson.M{"time": -1})

	if err := mgm.Coll(&NCovOverallInfo{}).SimpleFind(&infos, filter, findOptions); err != nil || len(infos) <= 0 {
		return nil, err
	}
	return infos[0], nil
}
