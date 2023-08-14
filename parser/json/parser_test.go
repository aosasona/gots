package jsonparser

import (
	"reflect"
	"testing"

	"github.com/aosasona/gots/v2/parser/tag"
)

type TestStruct struct {
	Name         string `json:"first_name"`
	OmitEmpty    string `json:",omitempty"`
	ShouldSkip   string `json:"-"`
	Invalid      string `json:"invalid,omitempty,invalid"`
	Default      string `json:""`
	Tagless      string
	Formed       string `json:"formed,omitempty"`
	AsStr        string `json:",string"`
	Dash         string `json:"-,"`
	WithNameOnly string `json:"name,"`
}

var testStruct = reflect.TypeOf(TestStruct{})

func TestJSONTagParser_Parse(t *testing.T) {
	ok := true

	nameField, ok := testStruct.FieldByName("Name")
	omitemptyField, ok := testStruct.FieldByName("OmitEmpty")
	skipField, ok := testStruct.FieldByName("ShouldSkip")
	invalidField, ok := testStruct.FieldByName("Invalid")
	defaultField, ok := testStruct.FieldByName("Default")
	taglessField, ok := testStruct.FieldByName("Tagless")
	propertlyFormedField, ok := testStruct.FieldByName("Formed")
	asStrField, ok := testStruct.FieldByName("AsStr")
	dashField, ok := testStruct.FieldByName("Dash")
	withNameOnlyField, ok := testStruct.FieldByName("WithNameOnly")

	if !ok {
		panic("field not found")
	}

	tests := []struct {
		Name     string
		Source   reflect.StructField
		Expected *tag.Tag
		WantErr  bool
	}{
		{
			Name:   "properly parse tag with custom name",
			Source: nameField,
			Expected: &tag.Tag{
				OriginalName: "Name",
				Name:         "first_name",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "should be optional",
			Source: omitemptyField,
			Expected: &tag.Tag{
				OriginalName: "OmitEmpty",
				Name:         "OmitEmpty",
				Skip:         false,
				Optional:     true,
			},
		},
		{
			Name:   "skip tag",
			Source: skipField,
			Expected: &tag.Tag{
				OriginalName: "ShouldSkip",
				Name:         "ShouldSkip",
				Skip:         true,
				Optional:     false,
			},
		},
		{
			Name:    "fail to parse invalid tag",
			Source:  invalidField,
			WantErr: true,
		},
		{
			Name:   "produce tag with default name",
			Source: defaultField,
			Expected: &tag.Tag{
				OriginalName: "Default",
				Name:         "Default",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly handle tagless fields",
			Source: taglessField,
			Expected: &tag.Tag{
				OriginalName: "Tagless",
				Name:         "Tagless",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly handle properly formed tag",
			Source: propertlyFormedField,
			Expected: &tag.Tag{
				OriginalName: "Formed",
				Name:         "formed",
				Skip:         false,
				Optional:     true,
			},
		},
		{
			Name:   "properly handle string tag",
			Source: asStrField,
			Expected: &tag.Tag{
				OriginalName: "AsStr",
				Name:         "AsStr",
				Skip:         false,
				Optional:     false,
				Type:         "string",
			},
		},
		{
			Name:   "properly handle dash as name tag",
			Source: dashField,
			Expected: &tag.Tag{
				OriginalName: "Dash",
				Name:         "-",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly handle tag with name only and empty second pair",
			Source: withNameOnlyField,
			Expected: &tag.Tag{
				OriginalName: "WithNameOnly",
				Name:         "name",
				Skip:         false,
				Optional:     false,
			},
		},
	}

	for _, test := range tests {
		got, err := Parse(test.Source)
		if (err != nil) != test.WantErr {
			t.Errorf("`%s`: unexpected error: %v", test.Name, err)
		}

		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf("`%s`: expected %+v, got %+v", test.Name, test.Expected, got)
		}
	}
}
