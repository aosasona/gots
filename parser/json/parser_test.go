package jsonparser

import (
	"reflect"
	"testing"

	"github.com/aosasona/gots/parser/tag"
)

type TestStruct struct {
	Name       string `json:"first_name"`
	OmitEmpty  string `json:",omitempty"`
	ShouldSkip string `json:"-"`
	Invalid    string `json:"invalid,omitempty,invalid"`
	Default    string `json:""`
	Tagless    string
	Formed     string `json:"formed,omitempty"`
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
	}

	for _, test := range tests {
		got, err := Parse(test.Source)
		if (err != nil) != test.WantErr {
			t.Errorf("failed to run case `%s`: unexpected error: %v", test.Name, err)
		}

		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf("failed to run case `%s`: expected %+v, got %+v", test.Name, test.Expected, got)
		}
	}
}
