package main

import (
	"log"
	"time"

	"github.com/aosasona/gots"
)

func main() {
	type Profession string

	type Person struct {
		FirstName  string     `json:"first_name"`
		LastName   string     `ts:"name:last_name"`
		DOB        string     `ts:"name:dob"`
		Profession Profession `ts:"name:job,optional:true"`
		CreatedAt  time.Time  // no tags
		DeletedAt  time.Time  `json:"-"`
		IsActive   bool       `ts:"name:is_active"`
		Ignored    []uint     `ts:"-"`
	}

	type Collection struct {
		CollectionName string   `ts:"name:name"`
		People         []Person `json:"whitelisted_users"` // an array of another struct
		Lead           Person
		Tags           []string `json:"omitempty" ts:"name:collection_tags"`
		AdminID        int      `json:"admin_id,omitempty"`
	}

	ts := gots.New(gots.Config{
		Enabled:           true,                   // you can use this to disable generation
		OutputFile:        "./example/index.d.ts", // this is where your generated file will be saved
		UseTypeForObjects: false,                  // if you want to use `type X = ...` instead of `interface X ...`
	})

	// registering a 'single' type
	// err := ts.Register(*new(Profession))
	// if err != nil {
	// 	log.Fatalf("error: %s\n", err.Error())
	// }

	// registering a normal interface
	// err = ts.Register(Person{})
	// if err != nil {
	// 	log.Fatalf("error: %s\n", err.Error())
	// }

	//registering a nested interface
	// err = ts.Register(Collection{})
	// if err != nil {
	// 	log.Fatalf("error: %s\n", err.Error())
	// }

	err := ts.Register(*new(Profession), Person{}, Collection{})
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
