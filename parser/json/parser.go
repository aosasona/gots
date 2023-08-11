package jsonparser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aosasona/gots/parser/tag"
)

type JSONTagParser struct{}

func (j *JSONTagParser) Parse(field reflect.StructField) (*tag.Tag, error) {
	tag := new(tag.Tag)
	tag.OriginalName = field.Name
	tag.Name = field.Name

	jsonTag := strings.TrimSpace(field.Tag.Get("json"))
	if jsonTag == "" {
		return tag, nil
	}

	if jsonTag == "-" {
		tag.Skip = true
		return tag, nil
	}

	// Parse with respect to the DEFAULT order, usually, the first field is the name and the second is to indicate if it should be optional
	values := strings.Split(jsonTag, ",")

	if len(values) > 2 {
		return nil, fmt.Errorf("expected at most 2 values, got %d", len(values))
	}

	if len(values) == 0 {
		return tag, nil
	}

	name := strings.TrimSpace(values[0])
	if name != "" {
		tag.Name = name
	}

	if len(values) > 1 {
		if strings.TrimSpace(values[1]) == "omitempty" {
			tag.Optional = true
		} else if strings.TrimSpace(values[1]) == "" {
			return nil, fmt.Errorf("expected omitempty or nothing, got %s", values[1])
		}
	}

	return tag, nil
}
