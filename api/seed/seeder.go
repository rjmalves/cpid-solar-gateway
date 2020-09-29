package seed

import (
	"math/rand"
	"time"

	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var inverters = []models.Inverter{
	{
		Serial:        "INVERTER1",
		Power:         500.0,
		Voltage:       127.0,
		Frequency:     60.0,
		Communication: true,
		Status:        true,
		Switch:        true,
	},
	{
		Serial:        "INVERTER2",
		Power:         750.0,
		Voltage:       127.0,
		Frequency:     60.0,
		Communication: true,
		Status:        true,
		Switch:        true,
	},
	{
		Serial:        "INVERTER3",
		Power:         250.0,
		Voltage:       127.0,
		Frequency:     60.0,
		Communication: true,
		Status:        true,
		Switch:        true,
	},
}

var telemetryData = []models.TelemetryData{}

// Load : loads the default seeded data into the DB
func Load(db *mongo.Database) error {
	// Adds the inverters to DB
	for _, inv := range inverters {
		_, err := inv.AddInverterToDB(db)
		if err != nil {
			return err
		}
	}
	// Generates a lot of TelemetryData
	for i := 0; i < 50; i++ {
		// Generates with random values and then attributes an inverter
		to := rand.Int63n(150)
		td := models.TelemetryData{
			Module:            "MODULE-X",
			LastTelemetryTime: time.Now().Unix() + to,
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
