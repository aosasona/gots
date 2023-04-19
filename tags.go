package gots

import (
	"reflect"
	"strings"
)

type parsedTag struct {
	Name     string
	Type     string
	Optional bool
	Skip     bool
	field    reflect.StructField
}

func parseTags(field reflect.StructField) parsedTag {
	parsed := newParsedTags(field)
	parsed.parseJSONTags().parseTSTag()
	return *parsed
}

func newParsedTags(field reflect.StructField) *parsedTag {
	return &parsedTag{field: field}
}

func (p *parsedTag) parseTSTag() *parsedTag {
	var tagFields []string

	tagFieldsMap := make(map[string]string)

	tag := strings.TrimSpace(p.field.Tag.Get("ts"))
	switch tag {
	case "":
		return p
	case "-":
		p.Skip = true
		return p
	}

	tagFields = strings.Split(tag, ",")
	if len(tagFields) == 0 {
		return p
	}

	for _, f := range tagFields {
		kv := strings.Split(f, ":")
		if len(kv) != 2 {
			continue
		}
		tagFieldsMap[kv[0]] = strings.TrimSpace(kv[1])
	}

	if name, ok := tagFieldsMap["name"]; ok {
		p.Name = name
	}

	if ty, ok := tagFieldsMap["type"]; ok {
		p.Type = ty
	}

	if optional, ok := tagFieldsMap["optional"]; ok {
		if optional == "true" || optional == "1" {
			p.Optional = true
		}
	}

	return p
}

func (p *parsedTag) parseJSONTags() *parsedTag {
	var tagFields []string

	tag := strings.TrimSpace(p.field.Tag.Get("json"))
	switch tag {
	case "":
		return p
	case "-":
		p.Skip = true
		return p
	}

	tagFields = strings.Split(tag, ",")
	if len(tagFields) == 0 {
		return p
	}

	switch tagFields[0] {
	case "omitempty":
		p.Optional = true
	default:
		p.Name = tagFields[0]
	}

	if len(tagFields) > 1 && tagFields[1] == "omitempty" {
		p.Optional = true
	}

	return p
}
