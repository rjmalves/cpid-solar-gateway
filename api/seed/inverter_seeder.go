package seed

import (
	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoadInverters : loads the default inverter data into the DB
func LoadInverters(db *mongo.Database) error {
	// Only seeds the DB if the collection is empty
	invs, err := models.ListInverters(db)
	if err != nil {
		return err
	}
	if len(invs) > 0 {
		return nil
	}
	// Creates the inverters
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
	// Adds the inverters to DB
	for _, inv := range inverters {
		_, err := inv.AddInverterToDB(db)
		if err != nil {
			return err
		}
	}
	return nil
}
