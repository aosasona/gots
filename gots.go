package gots

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aosasona/gots/config"
)

const (
	_STRING  = "string"
	_NUMBER  = "number"
	_BOOLEAN = "boolean"
	_ANY     = "any"
	_INVALID = "unknown"

	ErrNoSource = "no source specified"
)

type gots struct {
	config config.Config
	forks  []*gots
}

type parsedTag struct {
	Name     string
	Type     string
	Optional bool
	Skip     bool
}

func Init(c config.Config) *gots {
	if c.OutputFile == nil || *c.OutputFile == "" {
		c.OutputFile = config.String("types.ts")
	}

	return &gots{config: c}
}

// Fork takes the current gots instance and returns a new instance with the current config.
// If replace it true, it replaces the current config with the new config entirely, else it only replaces the non-nil values.
func (g *gots) Fork(c config.Config, replaceConfig bool) *gots {
	if replaceConfig {
		g.config = c
	} else {
		if c.Enabled != nil {
			g.config.Enabled = c.Enabled
		}

		if c.OutputFile != nil {
			g.config.OutputFile = c.OutputFile
		}

		if c.UseTypeForObjects != nil {
			g.config.UseTypeForObjects = c.UseTypeForObjects
		}

		if c.Case != "" {
			g.config.Case = c.Case
		}
	}

	fork := &gots{config: g.config}
	g.forks = append(g.forks, fork)

	return fork
}

func (g *gots) Register(sources ...any) error {
	if g.config.Enabled == nil || g.config.Enabled == config.Bool(false) {
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
