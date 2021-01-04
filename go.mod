module github.com/devcharmander/habit-tracker

go 1.14

replace github.com/devcharmander/habit-tracker/database => ../database

require (
	github.com/julienschmidt/httprouter v1.3.0
	go.mongodb.org/mongo-driver v1.4.4
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
