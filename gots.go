package gots

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	_STRING  = "string"
	_NUMBER  = "number"
	_BOOLEAN = "boolean"
	_ANY     = "any"
	_INVALID = "unknown"

	ERR_NO_SOURCE = "no source specified"
)

type gots struct {
	config Config
}

type Config struct {
	Enabled           bool
	OutputFile        string
	UseTypeForObjects bool
}

type parsedTag struct {
	Name     string
	Type     string
	Optional bool
	Skip     bool
}

func New(config Config) *gots {
	if config.OutputFile == "" {
		config.OutputFile = "index.d.ts"
	}

	return &gots{
		config,
	}
}

func (g *gots) Register(sources ...any) error {
	if !g.config.Enabled {
		return nil
	}

	if len(sources) == 0 {
		return errors.New(ERR_NO_SOURCE)
	}

	var output string
	for _, src := range sources {
		reflectType := reflect.TypeOf(src)

		if reflectType.Kind() == reflect.Struct {
			prefix := "export interface %s"
			if g.config.UseTypeForObjects {
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
