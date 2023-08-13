package gots

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aosasona/gots/config"
	"github.com/aosasona/gots/generator"
	"github.com/aosasona/gots/helper"
)

type Sources []any

type gots struct {
	config  config.Config
	sources Sources
}

var (
	// Deprecated: The `New` function has been replaced with `Init` but this is kept around for backwards compatibility with v1 (will be removed in a future release)
	New = Init

	// for convenience
	String = helper.String
	Bool   = helper.Bool

	ErrNoSources = errors.New("no sources provided")
)

const File_HEADER = `/**
* This file was generated by gots, do not edit it manually
* You can find the docs and source code at https://github.com/aosasona/gots
*/
	`

// for convenience
type Config = config.Config

func Init(c config.Config) *gots {
	if c.OutputFile == nil || *c.OutputFile == "" {
		c.OutputFile = helper.String("types.ts")
	}

	return &gots{config: c}
}

// Fork takes the current gots instance and returns a new instance with the current config.
// If replace it true, it replaces the current config with the new config entirely, else it only replaces the non-nil values.
func (g *gots) Fork(c config.Config, replaceConfig bool) *gots {
	fork := &gots{}

	if replaceConfig {
		fork.config = c
	} else {
		fork.config = g.config.Merge(c)
	}

	fork.sources = Sources{}

	return fork
}

func (g *gots) Count() int {
	return len(g.sources)
}

func (g *gots) AddSource(s any) {
	g.sources = append(g.sources, s)
}

func (g *gots) Commit(output string) error {
	defer func() { g.sources = Sources{} }()
	return nil
}

func (g *gots) Generate() (string, error) {
	var output string

	if len(g.sources) == 0 {
		return "", ErrNoSources
	}

	gn := generator.NewGenerator(generator.Opts{
		UseTypeForObjects: g.config.UseTypeForObjectsOrDefault(),
		ExpandStructs:     g.config.ExpandObjectTypesOrDefault(),
		PreferUnknown:     g.config.PreferUnknownOrDefault(),
	})

	for _, src := range g.sources {
		result := gn.Generate(src)
		output += result + "\n\n"
	}

	output = File_HEADER + "\n\n" + strings.TrimSpace(output)

	return output, nil
}

func (g *gots) Execute(log ...bool) error {
	output, err := g.Generate()
	if err != nil {
		return err
	}

	if len(log) > 0 && log[0] {
		fmt.Println(output)
	}

	return g.Commit(output)
}

// Calling Register will register the sources passed to it (doesn't replace the existing sources)
func (g *gots) Register(sources ...any) error {
	if !g.config.EnabledOrDefault() {
		return nil
	}

	if len(sources) == 0 {
		return ErrNoSources
	}

	g.sources = append(g.sources, sources...)

	return nil
}
