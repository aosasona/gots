package parser

import (
	"reflect"

	"go.trulyao.dev/gots/helper"
	gotsparser "go.trulyao.dev/gots/parser/gots"
	jsonparser "go.trulyao.dev/gots/parser/json"
	"go.trulyao.dev/gots/parser/tag"
)

func Parse(field reflect.StructField) (*tag.Tag, error) {
	tag := new(tag.Tag)

	tag.OriginalName = helper.WithDefaultString(tag.OriginalName, field.Name)
	tag.Name = helper.WithDefaultString(tag.Name, field.Name)

	// Parse the JSON struct tag first
	if _, err := jsonparser.Parse(field, tag); err != nil {
		return nil, err
	}

	// Parse the custom `gots` and `ts` struct tags to override the JSON struct tag if present
	if _, err := gotsparser.Parse(field, tag); err != nil {
		return nil, err
	}

	return tag, nil
}
