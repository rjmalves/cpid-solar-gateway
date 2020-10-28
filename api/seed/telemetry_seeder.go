package seed

import (
	"math/rand"
	"time"

	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoadTelemetryData : loads a lot of default telemetry data to the DB
func LoadTelemetryData(db *mongo.Database) error {
	// Only seeds the DB if the collection is empty
	tels, err := models.ListTelemetryData(db, bson.M{})
	if err != nil {
		return err
	}
	if len(tels) > 0 {
		return nil
	}
	// Generates a lot of TelemetryData
	n := 300
	tos := make([]int64, n)
	for i := range tos {
		tos[i] = int64(100 * i)
	}
	for i := range tos {
		// Generates with random values and then attributes an inverter
		td := models.TelemetryData{
			Module:            "MODULE-X",
			LastTelemetryTime: time.Now().Unix() + tos[i],
			OutputVoltage:     rand.Float64() * 280,
			InputVoltage:      rand.Float64() * 280,
			InputCurrent:      rand.Float64() * 10,
		}
		if i%3 == 0 {
			td.Serial = "INVERTER1"
		} else if i%3 == 1 {
			td.Serial = "INVERTER2"
		} else {
			td.Serial = "INVERTER3"
		}
		_, err := td.AddDataToDB(db)
		if err != nil {
			return err
		}
	}
	return nil
}
