package main

import (
	"fmt"
	"time"
)

func main() {
	Open()
	defer Close()

	// A Person struct consists of ID, Name, Age, Job.
	// var p ScienceNewsGroup
	// p.Data = []ScienceNews{}
	// p.Data = append(p.Data, ScienceNews{"Tidle1", "Summary1", "Url1"})
	// p.Data = append(p.Data, ScienceNews{"Title2", "Summary2", "Url2"})
	// p.save()
	// fmt.Println(p)

	var p ScienceNewsGroup
	err := p.Get(time.Now().String()[0:10])
	fmt.Println(p, err)

	// List("people")                    // each key/val in people bucket
	// ListPrefix("people", "10")        // ... with key prefix `20`
	// ListRange("people", "101", "101") // ... within range `101` to `103`
}
