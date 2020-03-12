package model

import (
	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type NCoVInfo struct {
	mgm.DefaultModel   `bson:",inline"`
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

func (*NCoVInfo) CollectionName() string {
	return "ncovinfos"
}

func (*NCoVInfo) Query(region, date string) ([]*NCoVInfo, error) {
	filter := bson.M{}
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
