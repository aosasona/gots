package generator

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aosasona/gots/v2/config"
	"github.com/aosasona/gots/v2/parser"
	"github.com/aosasona/gots/v2/parser/tag"
)

type TypeGenerator struct {
	expandStruct  bool // whether to expand struct types into object types like { foo: string } instead of the name
	preferUnknown bool // whether to prefer unknown over any
}

type TypeGeneratorOpts struct {
	ExpandStruct  bool
	PreferUnknown bool
}

type Property struct {
	Name         string
	Type         TSType
	OptionalChar string
}

var DefaultTypeGeneratorOpts = TypeGeneratorOpts{
	ExpandStruct:  false,
	PreferUnknown: true,
}

type TSType string

const (
	TypeString  TSType = "string"
	TypeNumber  TSType = "number"
	TypeBool    TSType = "boolean"
	TypeAny     TSType = "any"
	TypeUnknown TSType = "unknown"
	TypeObject  TSType = "object"

	// "Dynamic" types
	TypeArrayWithGenerics = "Array<%s>"
	TypeArray             = "%s[]"
	TypeRecord            = "Record<%s, %s>"

	// Nullable types
	TypeNullable = "%s | null"
	TypeOptional = "%s | undefined"
)

func NewTypeGenerator(opts TypeGeneratorOpts) *TypeGenerator {
	return &TypeGenerator{
		expandStruct:  opts.ExpandStruct,
		preferUnknown: opts.PreferUnknown,
	}
}

func (tg *TypeGenerator) GetFieldType(field reflect.StructField) TSType {
	return tg.getType(field)
}

func (tg *TypeGenerator) getType(field reflect.StructField) TSType {
	switch field.Type.Kind() {
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
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return TypeNumber
	case reflect.String:
		return TypeString
	case reflect.Bool:
		return TypeBool
	case reflect.Slice, reflect.Array:
		t := tg.getMappedType(field.Type.Elem())
		target := TypeArray
		if strings.Contains(string(t), "|") {
			target = TypeArrayWithGenerics
		}
		return TSType(fmt.Sprintf(target, t))
	case reflect.Interface:
		return TypeAny
	case reflect.Map:
		keyType := field.Type.Key()
		valueType := field.Type.Elem()
		return TSType(fmt.Sprintf(TypeRecord, tg.getMappedType(keyType), tg.getMappedType(valueType)))
	case reflect.Struct:
		if field.Type == reflect.TypeOf(time.Time{}) {
			return TypeNumber
		}
		return tg.toObjectType(field.Type)
	case reflect.Ptr:
		originalType := tg.getType(reflect.StructField{
			Name: field.Name,
			Type: field.Type.Elem(),
		})
		return TSType(fmt.Sprintf(TypeNullable, originalType))
	default:
		return tg.handleDefault(field)
	}
}

func (tg *TypeGenerator) getMappedType(src reflect.Type) TSType {
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
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return TypeNumber
	case reflect.String:
		return TypeString
	case reflect.Bool:
		return TypeBool
	case reflect.Interface:
		return TypeAny
	case reflect.Map:
		keyType := src.Key()
		valueType := src.Elem()
		return TSType(fmt.Sprintf(TypeRecord, tg.getMappedType(keyType), tg.getMappedType(valueType)))
	case reflect.Struct:
		if src == reflect.TypeOf(time.Time{}) {
			return TypeNumber
		}
		return tg.toObjectType(src)
	case reflect.Ptr, reflect.Slice, reflect.Array:
		if src.Elem().Kind() == reflect.Struct {
			return tg.toObjectType(src.Elem())
		} else if strings.Contains(src.String(), "*") {
			// this is a hack, but it works for now to see if it's a pointer to anything
			return tg.withNullable(tg.getMappedType(src.Elem()))
		}
		return tg.getMappedType(src.Elem())
	default:
		if tg.preferUnknown {
			return TypeUnknown
		}
		return TypeAny
	}
}

func (tg *TypeGenerator) withNullable(src TSType) TSType {
	return TSType(fmt.Sprintf(TypeNullable, src))
}

func (tg *TypeGenerator) toObjectType(src reflect.Type) TSType {
	// handle some known built-in types (like sql.Null*)
	switch src.String() {
	case "sql.NullString", "sql.NullByte":
		return tg.withNullable(TypeString)
	case "sql.NullInt16", "sql.NullInt32", "sql.NullInt64", "sql.NullFloat64", "sql.NullTime":
		return tg.withNullable(TypeNumber)
	case "sql.NullBool":
		return tg.withNullable(TypeBool)
	default:
		if tg.expandStruct {
			return tg.expandObjectType(src)
		}
		return TSType(src.Name())
	}
}

func (tg *TypeGenerator) expandObjectType(src reflect.Type) TSType {
	var fields []string

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)

		tag, err := parser.Parse(field)
		if err != nil {
			// if we are unable to parse the field, we should just skip it
			if config.Debug {
				fmt.Printf("unable to parse field: %s\n", field.Name)
			}
			continue
		}

		if tag.Skip {
			continue
		}

		p := makeProperty(field, tg, tag)

		if field.Anonymous {
			fields = append(fields, string(tg.expandObjectType(field.Type)))
		} else {
			fields = append(fields, fmt.Sprintf("%s%s%s: %s", TAB+TAB, p.Name, p.OptionalChar, p.Type))
		}
	}
	return TSType(fmt.Sprintf("{\n%s;\n%s}", strings.Join(fields, ";\n"), TAB))
}

func (tg *TypeGenerator) handleDefault(src reflect.StructField) TSType {
	if tg.preferUnknown {
		return TypeUnknown
	}
	return TypeAny
}

func makeProperty(field reflect.StructField, tg *TypeGenerator, tag *tag.Tag) Property {
	p := Property{
		Name: field.Name,
		Type: tg.getType(field),
	}

	if tag.Name != "" {
		p.Name = tag.Name
	}

	if tag.Optional {
		p.OptionalChar = "?"
	}

	if tag.Type != "" {
		p.Type = TSType(tag.Type)
	}

	return p
}
