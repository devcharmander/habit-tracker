package database

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//IClient - database client for habit tracker
type IClient interface {
	Add(data []*Habit) bool
	Update(data *Habit) bool
	Remove(data *Habit) bool
	Get(opts ...string) []*Habit
}

//Habit struct
type Habit struct {
	ID   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name" bson:"name"`
	Done bool          `json:"status" bson:"done"`
	Time time.Time     `json:"time" bson:"time"`
}
