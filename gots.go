package gots

import (
	"errors"

	"github.com/aosasona/gots/config"
	"github.com/aosasona/gots/helper"
)

type Sources []any

type gots struct {
	config  config.Config
	sources Sources
}

// Deprecated: The `New` function has been replaced with `Init` but this is kept around for backwards compatibility with v1 (will be removed in a future release)
var (
	New = Init

	ErrNoSources = errors.New("no sources provided")
)

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

func (g *gots) AddSource(s any) {
	g.sources = append(g.sources, s)
}

func (g *gots) Commit(output string) error {
	return nil
}

func (g *gots) Generate() (string, error) {
	if len(g.sources) == 0 {
		return "", ErrNoSources
	}
	return "", nil
}

func (g *gots) Execute() error {
	output, err := g.Generate()
	if err != nil {
		return err
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
