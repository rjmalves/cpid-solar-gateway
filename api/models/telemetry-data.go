package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var telemetryDataCollection = "telemetryData"

// TelemetryData : PV system state captured by the data acquisition service
type TelemetryData struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Serial            string             `bson:"serial" json:"serial"`
	Module            string             `bson:"module" json:"module"`
	LastTelemetryTime int64              `bson:"lastTelemetryTime" json:"lastTelemetryTime"`
	OutputVoltage     float64            `bson:"outputVoltage" json:"outputVoltage"`
	InputVoltage      float64            `bson:"inputVoltage" json:"inputVoltage"`
	InputCurrent      float64            `bson:"inputCurrent" json:"inputCurrent"`
}

// ListTelemetryData : reads telemetry data from DB using an filter
func ListTelemetryData(db *mongo.Database, filter bson.M) ([]*TelemetryData, error) {
	ctx := context.Background()
	cur, err := db.Collection(telemetryDataCollection).Find(ctx, filter)
	if err != nil {
		return []*TelemetryData{}, err
	}
	defer cur.Close(ctx)
	telemetry := []*TelemetryData{}
	for cur.Next(ctx) {
		var i TelemetryData
		if err := cur.Decode(&i); err != nil {
			return telemetry, err
		}
		telemetry = append(telemetry, &i)
	}
	return telemetry, nil
}

// AddDataToDB : adds a telemetry read to the DB
func (t *TelemetryData) AddDataToDB(db *mongo.Database) (primitive.ObjectID, error) {
	ctx := context.Background()
	res, err := db.Collection(telemetryDataCollection).InsertOne(ctx, t)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid, nil
}
