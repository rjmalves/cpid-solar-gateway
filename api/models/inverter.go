package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var inverterCollection = "inverters"

// Inverter : model of an inverter installed in the PV system
type Inverter struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Serial          string             `bson:"serial" json:"serial"`
	Power           float64            `bson:"power" json:"power"`
	Voltage         float64            `bson:"voltage" json:"voltage"`
	Frequency       float64            `bson:"frequency" json:"frequncy"`
	Communication   bool               `bson:"communication" json:"communication"`
	Status          bool               `bson:"status" json:"status"`
	Switch          bool               `bson:"switch" json:"switch"`
	EnergyToday     float64            `bson:"energyToday" json:"energyToday"`
	EnergyThisMonth float64            `bson:"energyThisMonth" json:"energyThisMonth"`
	EnergyThisYear  float64            `bson:"energyThisYear" json:"energyThisYear"`
	TotalEnergy     float64            `bson:"totalEnergy" json:"totalEnergy"`
}

// ListInverters : reads all the current inverters in the DB
func ListInverters(db *mongo.Database) ([]*Inverter, error) {
	ctx := context.Background()
	filter := bson.M{}
	cur, err := db.Collection(inverterCollection).Find(ctx, filter)
	if err != nil {
		return []*Inverter{}, err
	}
	defer cur.Close(ctx)
	inverters := []*Inverter{}
	for cur.Next(ctx) {
		var i Inverter
		if err := cur.Decode(&i); err != nil {
			return inverters, err
		}
		inverters = append(inverters, &i)
	}
	return inverters, nil
}

// ReadInverter : reads data from a specific inverter serial
func (i *Inverter) ReadInverter(db *mongo.Database) error {
	ctx := context.Background()
	filter := bson.M{
		"serial": i.Serial,
	}
	res := db.Collection(inverterCollection).FindOne(ctx, filter)
	if res.Err() != nil {
		return res.Err()
	}
	res.Decode(&i)
	return nil
}

// AddInverterToDB : adds info about a inverter to the DB
func (i *Inverter) AddInverterToDB(db *mongo.Database) (primitive.ObjectID, error) {
	ctx := context.Background()
	res, err := db.Collection(inverterCollection).InsertOne(ctx, i)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid, nil
}
