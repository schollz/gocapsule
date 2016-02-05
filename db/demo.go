package main

import (
	"fmt"
	"log"
)

func main() {
	Open()
	defer Close()

	// A Person struct consists of ID, Name, Age, Job.
	peeps := []*Person{
		{"100", "Bill Joy", "60", "Programmer"},
		{"101", "Peter Norvig", "58", "Programmer"},
		{"102", "Donald Knuth", "77", "Programmer"},
		{"103", "Jeff Dean", "47", "Programmer"},
		{"104", "Rob Pike", "59", "Programmer"},
		{"200", "Brian Kernighan", "73", "Programmer"},
		{"201", "Ken Thompson", "72", "Programmer"},
	}

	// Persist people in the database.
	for _, p := range peeps {
		p.save()
	}

	// Get a person from the database by their ID.
	for _, id := range []string{"100", "101"} {
		p, err := GetPerson(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(p.Name)
	}

	// List("people")                    // each key/val in people bucket
	// ListPrefix("people", "10")        // ... with key prefix `20`
	ListRange("people", "101", "101") // ... within range `101` to `103`
}
