package gots

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

const (
	_STRING  = "string"
	_NUMBER  = "number"
	_BOOLEAN = "boolean"
	_ANY     = "any"
	_INVALID = "unknown"

	ERR_NO_source = "no source specified"
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
}

func New(config Config) *gots {
	return &gots{
		config,
	}
}

func (g *gots) Register(sources ...any) error {
	if !g.config.Enabled {
		return nil
	}

	if len(sources) == 0 {
		return errors.New(ERR_NO_source)
	}

	var output string
	for _, src := range sources {
		reflectType := reflect.TypeOf(src)

		if reflectType.Kind() == reflect.Struct {
			prefix := "export interface %s"
			if g.config.UseTypeForObjects {
				prefix = "export type %s ="
			}

			output += fmt.Sprintf("%s %s;\n\n", fmt.Sprintf(prefix, reflectType.Name()), toObjectType(reflectType))
		} else {
			output += fmt.Sprintf("export %s\n\n", toSingleType(reflectType))
		}
	}

	log.Printf("output: %s", output)

	return nil
}

func toSingleType(src reflect.Type) string {
	return fmt.Sprintf("type %s = %s;", src.Name(), src.Kind())
}

func toObjectType(src reflect.Type) string {
	var fields []string
	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		parsedTags := parseFieldStructTag(field)
		mappedType := getMappedType(field.Type)

		fields = append(fields, makeTSInterfaceString(&field, mappedType, parsedTags))
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, "\n"))
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

	tagFields = strings.Split(tag, ",")
	if len(tagFields) == 0 {
		return result
	}

	for _, f := range tagFields {
		kv := strings.Split(f, ":")
		if len(kv) < 2 {
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
		if optional == "true" {
			result.Optional = true
		}
	}

	return result
}

func makeTSInterfaceString(field *reflect.StructField, mappedType string, override parsedTag) string {
	var (
		optionalChar string
		arrChar      string
	)
	name := field.Name

	if override.Name != "" {
		name = override.Name
	}

	if override.Type != "" {
		mappedType = override.Type
	}

	if override.Optional {
		optionalChar = "?"
	}

	if field.Type.Kind() == reflect.Slice {
		arrChar = "[]"
	}

	return fmt.Sprintf("%s%s: %s%s;", name, optionalChar, mappedType, arrChar)
}

func getMappedType(src reflect.Type) string {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return _NUMBER
	case reflect.String:
		return _STRING
	case reflect.Bool:
		return _BOOLEAN
	case reflect.Interface:
		return _ANY
	case reflect.Struct:
		if src == reflect.TypeOf(time.Time{}) {
			return _STRING
		}
		return toObjectType(src)
	case reflect.Ptr, reflect.Slice, reflect.Array:
		return getMappedType(src.Elem())
	default:
		return _INVALID
	}
}
