package db

import (
	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDB(uri, name string) error {
	err := mgm.SetDefaultConfig(nil, name, options.Client().ApplyURI(uri))
	return err
}
