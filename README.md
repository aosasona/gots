# GoTS

## So, what is GoTS

GoTS allows you to generate TypeScript types from your selected Golang Structs in the Golang code itself. 

## Why Would I do this?

> Let me paint you a picture

You have embedded React.js or Astro in your Golang application (or you have them together in a Monorepo) but now you have to define Typescript types for your Golang's API responses that you already have structs for.
Fine, you may feel okay with doing that, what happens when you change one of those struct types? Now you need to update the matching TS type so you don't shoot yourself in the foot. 
But you're human, you could easily forget and that's not good. That's where this comes in it generates the types during run-time which means it will always be up-to-date especially if you use something like air for hot-reloading.

## How do I use this?

It's very easy to use, here's how:

```go
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
		collectionName string   `ts:"name:name"`
		people         []Person // an array of another struct
		lead           Person
	}

	ts := gots.New(gots.Config{
		Enabled:           true,           // you can use this to disable generation
		OutputFile:        "./index.d.ts", // this is where your generated file will be saved
		UseTypeForObjects: false,          // if you want to use `type X = ...` instead of `interface X ...`
	})

	// registering a 'single' type
	err := ts.Register(*new(Profession), Person{}, Collection{})
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
```

## Contribution

PRs and issues are welcome. You can find me on Twiter at [@trulyao](https://twitter.com/trulyao)
