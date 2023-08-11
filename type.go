package gots

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func toSingleType(src reflect.Type) string {
	kind := getMappedType(src)
	return fmt.Sprintf("type %s = %s;", src.Name(), kind)
}

func toObjectType(src reflect.Type) string {
	var fields []string
	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		parsedTags := parseFieldStructTag(field)

		if parsedTags.Skip {
			continue
		}

		mappedType := getMappedType(field.Type)

		fields = append(fields, makeTSInterfaceString(&field, mappedType, parsedTags))
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, "\n"))
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

	return fmt.Sprintf("\t%s%s: %s%s;", name, optionalChar, mappedType, arrChar)
}

func getMappedType(src reflect.Type) string {
	switch src.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
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
