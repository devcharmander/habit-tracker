package database

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

const layoutISO = "2006-01-02"

//MongoClient struct
type MongoClient struct {
	Database       string
	CollectionName string
	Client         *mongo.Client
	Disconnect     func()
	Context        context.Context
	Collection     *mongo.Collection
}

//NewClient creates a database client with url and type specified in the params
func NewClient(clientType, url string) IClient {
	switch clientType {
	case "mongo":
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			log.Panicf("Could not connect to the url %s", url)
		}
		disconnect := func() {
			if err := client.Disconnect(ctx); err != nil {
				cancel()
				log.Panicf("Error disconnecting client.Err %v", err)
			}
			cancel()
		}

		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			log.Panicf("Could not ping the database. Error %v", err)
		}
		collection := client.Database("timetable").Collection("habits")
		return &MongoClient{
			Client:     client,
			Disconnect: disconnect,
			Context:    ctx,
			Collection: collection,
		}
	default:
		return nil
	}
}

//Add adds a Habit into database
func (m *MongoClient) Add(data []*Habit) bool {
	defer m.Disconnect()
	for _, record := range data {
		_, err := m.Collection.InsertOne(m.Context, record)
		if err != nil {
			log.Printf("Could not insert record. Errror %v", err)
			return false
		}
	}
	return true
}

//Update updates the done value in the database for a Habit
func (m *MongoClient) Update(data *Habit) bool {
	defer m.Disconnect()
	res, err := m.Collection.UpdateOne(
		m.Context,
		bson.M{"_id": data.ID},
		bson.D{{"$set", bson.D{{"done", data.Done}}}},
	)
	if err != nil {
		log.Printf("Could not update data. Error: %v", err)
		return false
	}
	if res.UpsertedID != nil {
		return true
	}
	return false
}

//Remove removes a Habit
func (m *MongoClient) Remove(data *Habit) bool {
	defer m.Disconnect()
	res, err := m.Collection.DeleteMany(m.Context, bson.M{"name": data.Name})
	if err != nil {
		log.Printf("Could not delete %s. Error: %v", data.Name, err)
		return false
	}
	if res.DeletedCount != 0 {
		return true
	}
	return false
}

//Get gets the data for the current day by default
//opts are slice of strings denoting startdate and enddate e.g. Get(01-02-2020,07-02-2020)
func (m *MongoClient) Get(opts ...string) []*Habit {
	defer m.Disconnect()
	habits := make([]*Habit, 0)
	var filter interface{}
	if len(opts) == 2 {
		startDate, err := time.Parse(layoutISO, opts[0])
		if err != nil {
			log.Println("Error parsing the time provided in options. Error", err)
			return habits
		}
		endDate, err := time.Parse(layoutISO, opts[1])
		if err != nil {
			log.Println("Error parsing the time provided in options. Error", err)
			return habits
		}
		filter = bson.M{"time": bson.M{"$gte": startDate, "$lte": endDate}}
	} else {
		date, err := time.Parse(layoutISO, strings.Split(time.Now().String(), " ")[0])
		if err != nil {
			log.Print("Error parsing date. Error", err)
			return habits
		}
		//fromDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		filter = bson.M{"time": bson.M{"$eq": date}}
	}
	cur, err := m.Collection.Find(m.Context, filter)
	if err != nil {
		log.Printf("could not find records for the filter %v. Error %v", filter, err)
		return habits
	}
	for cur.Next(m.Context) {
		habit := &Habit{}
		err := cur.Decode(habit)
		if err != nil {
			log.Printf("Could not decode the data into struct 'Habit'. Error %v", err)
		}
		habits = append(habits, habit)
	}
	return habits
}
