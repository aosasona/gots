package gots

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aosasona/gots/config"
	"github.com/aosasona/gots/helper"
)

const (
	_STRING  = "string"
	_NUMBER  = "number"
	_BOOLEAN = "boolean"
	_RECORD  = "Record<string, any>"
	_ANY     = "any"
	_INVALID = "unknown"

	ErrNoSource = "no source specified"
)

type Sources []any

type gots struct {
	config  config.Config
	output  []byte
	sources Sources
	forks   []*gots
}

type parsedTag struct {
	Name     string
	Type     string
	Optional bool
	Skip     bool
}

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
	fork.forks = []*gots{}
	fork.output = []byte{}

	g.forks = append(g.forks, fork)

	return fork
}

func (g *gots) AddSource(s any) {
	g.sources = append(g.sources, s)
}

func (g *gots) Commit(output string) error {
	return nil
}

func (g *gots) Register(sources ...any) error {
	if g.config.Enabled == nil || !*g.config.Enabled {
		return nil
	}

	if len(sources) == 0 {
		return errors.New(ErrNoSource)
	}

	var output string
	for _, src := range sources {
		reflectType := reflect.TypeOf(src)

		if reflectType.Kind() == reflect.Struct {
			prefix := "export interface %s"
			if g.config.UseTypeForObjectsOrDefault() {
				prefix = "export type %s ="
			}

			output += fmt.Sprintf("%s %s\n\n", fmt.Sprintf(prefix, reflectType.Name()), toObjectType(reflectType))
		} else {
			output += fmt.Sprintf("export %s\n\n", toSingleType(reflectType))
		}
	}

	err := g.exportToFile(output)

	// TODO: call commit here too

	return err
}

func parseFieldStructTag(field reflect.StructField) parsedTag {
	var (
		result    parsedTag
		tagFields []string
	)

	tagFieldsMap := make(map[string]string)

	tag := field.Tag.Get("ts")
	if tag == "" {
		return result
	}

	if tag == "-" {
		result.Skip = true
		return result
	}

	tagFields = strings.Split(tag, ",")
	if len(tagFields) == 0 {
		return result
	}

	for _, f := range tagFields {
		kv := strings.Split(f, ":")
		if len(kv) != 2 {
			continue
		}
		tagFieldsMap[kv[0]] = strings.TrimSpace(kv[1])
	}

	if name, ok := tagFieldsMap["name"]; ok {
		result.Name = name
	}

	if ty, ok := tagFieldsMap["type"]; ok {
		result.Type = ty
	}

	if optional, ok := tagFieldsMap["optional"]; ok {
		if optional == "true" || optional == "1" {
			result.Optional = true
		}
	}

	return result
}
