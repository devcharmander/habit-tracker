package database

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	habits := make([]*Habit, 0)
	date, err := time.Parse(layoutISO, strings.Split(time.Now().String(), " ")[0])
	if err != nil {
		print(err)
		t.Fatal()
	}
	for i := 0; i <= 10; i++ {
		for j := 0; j <= 7; j++ {
			habit := &Habit{
				Name: fmt.Sprintf("Test Habit %d", j),
				Time: date.AddDate(0, 0, i),
				Done: true,
			}
			habits = append(habits, habit)
		}
	}
	client := NewClient("mongodb://localhost:27017")
	res := client.Add(habits)
	if !res {
		t.Fatalf("Could not add one or more records")
	}
}

func TestGet(t *testing.T) {
	client := NewClient("mongodb://localhost:27017")
	habits := client.Get()
	if len(habits) == 0 {
		t.Fatal("Error getting habits")
	}
}

func TestGetOnRange(t *testing.T) {
	client := NewClient("mongodb://localhost:27017")
	habits := client.Get("2020-12-29", "2020-12-30")
	if len(habits) != 2 {
		log.Fatalf("Expected habits: 2. Actual Habits: %d", len(habits))
	}
	printHabits(habits)
}

func TestUpdate(t *testing.T) {
	client := NewClient("mongodb://localhost:27017")
	habits := client.Get()
	habit := habits[0]
	habit.Name = "Modified name"
	habit.Done = false
	client = NewClient("mongodb://localhost:27017")
	if !client.Update(habit) {
		log.Fatal("Expected true actual false")
	}
}

func TestDelete(t *testing.T) {
	client := NewClient("mongodb://localhost:27017")
	habits := client.Get()
	if len(habits) != 1 {
		log.Fatal("Could not get habit for the current date")
	}
	habit := habits[0]
	client = NewClient("mongodb://localhost:27017")
	if !client.Remove(habit) {
		log.Printf("could not delete data for %s", habit.Name)
	}
}

func printHabits(habits []*Habit) {
	for _, v := range habits {
		fmt.Printf("ID: %s\t\tName: %s\t\t time: %v\n", v.ID.Hex(), v.Name, v.Time)
	}
}
