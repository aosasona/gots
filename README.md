![gots](./assets/gots.png)

## So, what is gots?

gots allows you to generate TypeScript types from your selected Go types (int, string, struct etc) in the code itself. 

## Why would I do this?

> Let me paint you a picture

You have embedded React.js or Astro in your Go application (or you have them together in a monorepo) but now you have to define Typescript types for your Go's API responses (or other things) that you already have structs for.

Fine, you may feel okay with doing that, what happens when you change one of those types in the Go code? Now you need to update the matching TS type so you don't shoot yourself in the foot. 

But you're human, you could easily forget and that's not great now, is it? That's where this package comes in, it generates the types during run-time which means it will always be up-to-date especially if you use something like air for hot-reloading.

View generated example [here](./example/index.d.ts)

# Installation
Just paste this in your terminal (I promise it's safe):
```bash
go get github.com/aosasona/gots
```

# Usage

It's fairly easy to use, here's how:

```go
package main

import (
	"log"
	"os"
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

	// registering multiple types at once
	err := ts.Register(*new(Profession), Person{}, Collection{})
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
```

You can pass in the following override values via struct field tags:
- name (string)
- type (string)
- optional (only `true` or `1` or it is ignored)

These give you more control over what types end up being generated. You don't need to specify these, they optional, if they are not specified the default values are inferred from the types themselves.

It is safer to enable gots in development only, you can do this however way you want in your application. For example:

```go
...
ts := gots.New(gots.Config{
	Enabled: os.Getenv("ENV") == "development",
})
...
```


## Contribution

PRs and issues are welcome. You can find me on Twiter at [@trulyao](https://twitter.com/trulyao)
