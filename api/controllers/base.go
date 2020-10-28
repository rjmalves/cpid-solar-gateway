package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Server : the base elements that make the service
type Server struct {
	DB     *mongo.Database
	Router *gin.Engine
}

// Initialize : prepares the service to launch
func (s *Server) Initialize(DBHost, DBPort, DBUser, DBPassword, DBDatabase, mode string) error {
	// Connects with the database
	mongoURI := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", DBUser, DBPassword, DBHost, DBPort, DBDatabase)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	s.DB = client.Database(DBDatabase)
	// Checks if first connection for making the collection setups
	colls, err := s.DB.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}
	inverterCollFound := false
	telemetryCollFound := false
	for _, c := range colls {
		if c == "inverters" {
			inverterCollFound = true
		} else if c == "telemetryData" {
			telemetryCollFound = true
		}
	}
	shouldSetup := !(inverterCollFound && telemetryCollFound)
	if shouldSetup {
		if err := s.DBSetup(ctx); err != nil {
			return err
		}
	}
	// Checks if the mode is valid for running the server
	validModes := []string{"debug", "release"}
	valid := false
	for _, m := range validModes {
		if mode == m {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("Service execution mode not valid")
	}
	gin.SetMode(mode)
	// Creates the HTTP router
	s.Router = gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "HEAD", "PUT", "POST", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Access-Control-Allow-Headers",
		"Origin",
		"Authorization",
		"Accept",
		"X-Requested-With",
		"Content-Type",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers"}
	s.Router.Use(cors.New(config))
	s.initializeRoutes()

	return nil
}

// Terminate : closes connections and ends the service
func (s *Server) Terminate() error {
	// Disconnects from DB
	if err := s.DB.Client().Disconnect(context.Background()); err != nil {
		return err
	}
	return nil
}

// Run : runs the service and recovers errors
func (s *Server) Run(servicePort string) error {
	defer s.Terminate()
	// Hosts the API
	go func() {
		addr := ":" + servicePort
		for {
			if err := s.Router.Run(addr); err != nil {
				log.Printf("Error serving the API: %v\n", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	// Exits on SIGINT
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	return nil
}

// DBSetup : setups the DB collections in the first launch
func (s *Server) DBSetup(ctx context.Context) error {
	// Creates collections with existence rules
	opts := options.CreateCollection()
	opts.SetCapped(true)
	opts.SetSizeInBytes(1e11)
	if err := s.DB.CreateCollection(ctx, "inverters", opts); err != nil {
		return err
	}
	if err := s.DB.CreateCollection(ctx, "telemetryData", opts); err != nil {
		return err
	}
	iCol := s.DB.Collection("inverters")
	tCol := s.DB.Collection("telemetryData")
	// Creates unique indexes
	iMod := mongo.IndexModel{
		Keys: bson.M{
			"serial": -1,
		},
		Options: options.Index().SetUnique(true),
	}
	tMod := mongo.IndexModel{
		Keys: bson.M{
			"serial":            -1,
			"lastTelemetryTime": -1,
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := iCol.Indexes().CreateOne(ctx, iMod); err != nil {
		return err
	}
	if _, err := tCol.Indexes().CreateOne(ctx, tMod); err != nil {
		return err
	}
	return nil
}
