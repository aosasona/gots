package main

import (
	"log"
	"time"

	"github.com/aosasona/gots"
)

func main() {
	type Profession string

	type Person struct {
		firstName  string `ts:"name:first_name"`
		lastName   string `ts:"name:last_name"`
		dob        string
		profession Profession `ts:"name:job,optional:true"`
		createdAt  time.Time
		isActive   bool `ts:"name:is_active"`
	}

	type Collection struct {
		collectionName string `ts:"name:name"`
		people         []Person
	}

	ts := gots.New(gots.Config{
		Enabled: true,
	})

	// registering a 'single' type
	err := ts.Register(*new(Profession))
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}

	// registering a normal interface
	err = ts.Register(Person{})
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}

	// registering a nested interface
	err = ts.Register(Collection{})
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
