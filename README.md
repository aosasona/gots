> [!WARNING]
> This package has been renamed and because of that, development has moved [here](https://github.com/aosasona/mirror). This repo has been archived to prevent breaking your current usage.

![gots](./assets/gots.png)

## So, what is gots?

gots allows you to generate usable TypeScript types from your selected Go types (int, string, struct etc) in the code itself.

## Why would I do this?

> Let me paint you a picture

You have embedded React.js or Astro in your Go application (or you have them together in a monorepo) but now you have to define Typescript types for your Go's API responses (or other things) that you already have types for in your Go code.

Fine, you may feel okay with doing that, what happens when you change one of those types in the Go code? Now you need to update the matching TS type so you don't shoot yourself in the foot.

But you're human, you could easily forget and that's not great now, is it? That's where this package comes in, it generates the types during run-time which means it will always be up-to-date especially if you use something like air for hot-reloading.

> plus, I just enjoy building stuff like this

View generated example [here](./examples/example.ts)

## You should know...

The generated types may not always match what you expect (especially in the cases of embedded structs) and might just be an `any` or `unknown`, to be more specific, it is advised to use the type property in the `gots` or `ts` struct tag to specify the type yourself
Gots is not designed or built to be or ever be 100% accurate, just enough to have you setup and ready to communicate with your Go service/app/API _safely_ in Typescript, knowing a large part of what to send and expect back.

# Installation

Just paste this in your terminal (I promise it's safe):

```bash
go get -u github.com/aosasona/gots/v2
```

# Usage

Not an exceptional documentation but this should help you get started

```go
package main

import (
	"fmt"
	"time"

	"github.com/aosasona/gots"
	"github.com/aosasona/gots/config"
)

type Language string

type Tags map[string]string

type Person struct {
	FName     string         `gots:"name:first_name"`
	LName     string         `gots:"name:last_name"`
	Age       int            `gots:"name:age"`
	Languages []Language     `gots:"name:languages"`
	Grades    map[string]int `gots:"name:grades,optional:1"`
	Tags      Tags           `gots:"name:tags"`
	CreatedAt time.Time      `gots:"name:created_at"`
	UpdatedAt *time.Time     `gots:"name:updated_at"`
	DeletedAt *time.Time     `gots:"name:deleted_at"`
}

func main() {
	gt := gots.Init(config.Config{
		Enabled:           gots.Bool(true),
		OutputFile:        gots.String("./examples/example.ts"),
		UseTypeForObjects: gots.Bool(true),
		ExpandObjectTypes: gots.Bool(true),
	})

	// ===> Individually
	gt.AddSource(*new(Language))
	gt.AddSource(Tags{})
	gt.AddSource(Person{})

	out, err := gt.Generate()
	if err != nil {
		fmt.Println(err)
		return
	}

	// save to file
	err = gt.Commit(out)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ===> As a group
	gt.Register(*new(Language), Tags{}, Person{})

	// generate types and save to the file
	err := gt.Execute(true)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

It is safer to enable gots in development only, you can do this however way you want in your application. For example:

```go
...
ts := gots.Init(gots.Config{
	Enabled: os.Getenv("ENV") == "development",
})
...
```

## Tags

You can configure the generated types using struct tags; the `json` tag, the `gots` tag or the legacy `ts` struct tag. You can pass in the following override values via struct field tags:

- name (string)
- type (string)
- optional (only `true` or `1` or it is ignored)
- skip (only `true` or `1`, but can also simply be written like this: `gots:"-"`)

#### Example

```go
type Ex struct {
	ID	string `json:"user_id,omitempty" gots:"type:Uppercase<string>"`
	Name string `gots:"name:fname,optional:true"`
}
```

This will translate into:

```typescript
export interface Ex {
	user_id?: Uppercase<string>;
	fname?: string;
}
```

These give you more control over what types end up being generated. You don't need to specify these, they are optional, if they are not specified, the default values are inferred from the types themselves.

## Contribution

PRs and issues are welcome :)

## Development

- To run the example:

```sh
just example
```

- To run the tests:

```sh
just test
```
