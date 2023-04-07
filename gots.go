package gots

import (
	"errors"
	"fmt"
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
	Enabled   bool
	OutputDir string
	Overwrite bool
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

func (*gots) Register(sources ...any) error {
	if len(sources) == 0 {
		return errors.New(ERR_NO_source)
	}

	var output string

	for _, src := range sources {
		reflectType := reflect.TypeOf(src)
		if reflectType.Kind() == reflect.Struct {
			output = toObjectType(reflectType)
		} else {
			output = toSingleType(reflectType)
		}
	}

	fmt.Printf("output: %s", output)

	return nil
}

func toSingleType(src reflect.Type) string {
	return fmt.Sprintf("type %s = %s;", src.Name(), src.Kind())
}

func toObjectType(src reflect.Type) string {
	var fields []string
	for i := 0; i < src.NumField(); i++ {
	}

	return fmt.Sprintf("interface %s {\n%s\n}", src.Name(), strings.Join(fields, "\n"))
}

func parseStructTag(field reflect.StructField) parsedTag {
	var result parsedTag

	tag := field.Tag.Get("ts")
	if tag == "" {
		return result
	}

	return result
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
