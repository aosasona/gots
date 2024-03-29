package generator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aosasona/gots/v2/config"
	"github.com/aosasona/gots/v2/parser"
)

// acts as a tab "character" for formatting purposes in the generated code (4 spaces)
const TAB = "    "

var tg *TypeGenerator

type Generator struct {
	useTypeForObjects bool
}

type Opts struct {
	UseTypeForObjects bool
	ExpandStructs     bool // whether to expand struct types into object types like { foo: string } instead of the name
	PreferUnknown     bool // whether to prefer unknown over any
}

func NewGenerator(opts Opts) *Generator {
	tg = NewTypeGenerator(TypeGeneratorOpts{
		ExpandStruct:  opts.ExpandStructs,
		PreferUnknown: opts.PreferUnknown,
	})

	return &Generator{
		useTypeForObjects: opts.UseTypeForObjects,
	}
}

func (g *Generator) Generate(src any) string {
	var (
		srcType = reflect.TypeOf(src)
		result  string
	)

	if srcType.Kind() == reflect.Struct {
		result = "export interface %s "
		if g.useTypeForObjects {
			result = "export type %s = "
		}

		result += g.generateObjectType(srcType)
	} else {
		result = "export type %s = "
		result += string(tg.GetFieldType(reflect.StructField{
			Name: srcType.Name(),
			Type: srcType,
		}))
	}

	return fmt.Sprintf(result, srcType.Name())
}

func (g *Generator) generateObjectType(src reflect.Type) string {
	var result string

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)

		tag, err := parser.Parse(field)
		if err != nil {
			if config.Debug {
				fmt.Printf("Error parsing field: %s\n", err.Error())
			}
			continue
		}

		if !field.IsExported() || tag.Skip {
			continue
		}

		p := makeProperty(field, tg, tag)

		result += fmt.Sprintf("%s%s%s: %s;\n", TAB, p.Name, p.OptionalChar, p.Type)
	}

	return fmt.Sprintf("{\n%s%s\n}", TAB, strings.TrimSpace(result))
}
