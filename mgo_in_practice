// This program provides a sample application for using MongoDB with
// the mgo driver.
package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"sync"
	"time"
)

const (
	MongoDBHosts = "ds035428.mongolab.com:35428"
	AuthDatabase = "goinggo"
	AuthUserName = "guest"
	AuthPassword = "welcome"
	TestDatabase = "goinggo"
)

type (
	// BuoyCondition contains information for an individual station.
	BuoyCondition struct {
		WindSpeed     float64 `bson:"wind_speed_milehour"`
		WindDirection int     `bson:"wind_direction_degnorth"`
		WindGust      float64 `bson:"gust_wind_speed_milehour"`
	}

	// BuoyLocation contains the buoy's location.
	BuoyLocation struct {
		Type        string    `bson:"type"`
		Coordinates []float64 `bson:"coordinates"`
	}

	// BuoyStation contains information for an individual station.
	BuoyStation struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		StationId string        `bson:"station_id"`
		Name      string        `bson:"name"`
		LocDesc   string        `bson:"location_desc"`
		Condition BuoyCondition `bson:"condition"`
		Location  BuoyLocation  `bson:"location"`
	}
)

// main is the entry point for the application.
func main() {
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	mongoSession.SetMode(mgo.Monotonic, true)

	// Create a wait group to manage the goroutines.
	var waitGroup sync.WaitGroup

	// Perform 10 concurrent queries against the database.
	waitGroup.Add(10)
	for query := 0; query < 10; query++ {
		go RunQuery(query, &waitGroup, mongoSession)
	}

	// Wait for all the queries to complete.
	waitGroup.Wait()
	log.Println("All Queries Completed")
}

// RunQuery is a function that is launched as a goroutine to perform
// the MongoDB work.
func RunQuery(query int, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) {
	// Decrement the wait group count so the program knows this
	// has been completed once the goroutine exits.
	defer waitGroup.Done()

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(TestDatabase).C("buoy_stations")

	log.Printf("RunQuery : %d : Executing\n", query)

	// Retrieve the list of stations.
	var buoyStations []BuoyStation
	err := collection.Find(nil).All(&buoyStations)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

	log.Printf("RunQuery : %d : Count[%d]\n", query, len(buoyStations))
}
