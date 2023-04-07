package main

import (
	"log"
	"time"

	"github.com/aosasona/gots"
)

func main() {
	type Profession string

	type Person struct {
		firstName  string
		lastName   string     `ts:"name:last_name"`
		age        int        `ts:"type:string"`
		profession Profession `ts:"name:job"`
		createdAt  time.Time  `ts:"optional:true"`
	}

	type NestedStruct struct {
		collectionName string
		people         []Person
	}

	ts := gots.New(gots.Config{
		Enabled: true,
	})

	// registering a 'single' type
	err := ts.Register(*new(Profession))
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}
