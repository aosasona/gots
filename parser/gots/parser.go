package gotsparser

import (
	"reflect"
	"strings"

	"github.com/aosasona/gots/helper"
	"github.com/aosasona/gots/parser/tag"
)

func withDefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

func Parse(field reflect.StructField, targetTag ...*tag.Tag) (*tag.Tag, error) {
	var (
		tag  = new(tag.Tag)
		opts = make(map[string]string)
	)

	if len(targetTag) > 0 {
		tag = targetTag[0]
	} else {
		tag.OriginalName = helper.WithDefaultString(tag.OriginalName, field.Name)
		tag.Name = helper.WithDefaultString(tag.Name, field.Name)
	}

	gotsTag := strings.TrimSpace(field.Tag.Get("gots"))

	// check if there is a `ts` tag in the field (for backwards compatibility)
	if gotsTag == "" {
		gotsTag = strings.TrimSpace(field.Tag.Get("ts"))

		if gotsTag == "" {
			return tag, nil
		}
	}

	if gotsTag == "-" {
		tag.Skip = true
		return tag, nil
	}

	tagFields := strings.Split(gotsTag, ",")
	if len(tagFields) == 0 {
		return tag, nil
	}

	// split the props into key-value pairs
	for _, f := range tagFields {
		kv := strings.Split(f, ":")
		if len(kv) != 2 {
			continue
		}
		opts[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	if skip, ok := opts["skip"]; ok {
		if strings.TrimSpace(skip) == "true" || strings.TrimSpace(skip) == "1" {
			tag.Skip = true

			// if the tag is set to skip, then we don't need to parse the rest of the tag
			return tag, nil
		}
	}

	if name, ok := opts["name"]; ok {
		tag.Name = withDefaultString(strings.TrimSpace(name), field.Name)
	}

	if optional, ok := opts["optional"]; ok {
		// check if the optional tag is set to true or 1 (for backwards compatibility)
		if strings.TrimSpace(optional) == "true" || strings.TrimSpace(optional) == "1" {
			tag.Optional = true
		}
	}

	if overrideType, ok := opts["type"]; ok {
		tag.Type = strings.TrimSpace(overrideType)
	}

	return tag, nil
}
