package generator

import (
	"database/sql"
	"reflect"
	"testing"
	"time"
)

type Tag struct {
	Name string
	Type string
}

type Address struct {
	Street   string `ts:"name:street,optional:1"`
	City     string `json:"city,omitempty"`
	Postcode string `gots:"name:post_code,type:Capitalize<string>"`
}

type SocialMedia struct {
	Instagram *string  `json:"instagram_handle"`
	Twitter   string   `json:"twitter_handle"`
	Blogs     []string `gots:"name:blogs"`
	HNItems   []*int   `json:"hackernews_items"`
}

type Person struct {
	Name          string
	Age           int
	EmailVerified bool
	Address       Address
	Contact       []string
	Scores        []int
	Connections   []*uint
	TagsMap       map[string]Tag
	Tags          []Tag
	Milestones    map[int]time.Time
	Parent        sql.NullInt64
	CreatedAt     time.Time
	DeletedAt     *time.Time
	SocialMedia
}

func Test_GetTypeUsingDefaultOpts(t *testing.T) {
	var name string
	p := Person{}

	tests := []struct {
		Name     string
		Source   reflect.StructField
		Expected TSType
	}{
		{
			Name: "string",
			Source: reflect.StructField{
				Name: "Name",
				Type: reflect.TypeOf(p.Name),
			},
			Expected: TypeString,
		},
		{
			Name: "number",
			Source: reflect.StructField{
				Name: "Age",
				Type: reflect.TypeOf(p.Age),
			},
			Expected: TypeNumber,
		},
		{
			Name: "bool",
			Source: reflect.StructField{
				Name: "EmailVerified",
				Type: reflect.TypeOf(p.EmailVerified),
			},
			Expected: TypeBool,
		},
		{
			Name: "struct",
			Source: reflect.StructField{
				Name: "Address",
				Type: reflect.TypeOf(p.Address),
			},
			Expected: TSType("Address"),
		},
		{
			Name: "int slice",
			Source: reflect.StructField{
				Name: "Contact",
				Type: reflect.TypeOf(p.Contact),
			},
			Expected: TSType("string[]"),
		},
		{
			Name: "int slice",
			Source: reflect.StructField{
				Name: "Scores",
				Type: reflect.TypeOf(p.Scores),
			},
			Expected: TSType("number[]"),
		},
		{
			Name: "struct slice",
			Source: reflect.StructField{
				Name: "Tags",
				Type: reflect.TypeOf(p.Tags),
			},
			Expected: TSType("Tag[]"),
		},
		{
			Name: "pointer slice",
			Source: reflect.StructField{
				Name: "Connections",
				Type: reflect.TypeOf(p.Connections),
			},
			Expected: TSType("Array<number | null>"),
		},
		{
			Name: "map",
			Source: reflect.StructField{
				Name: "TagsMap",
				Type: reflect.TypeOf(p.TagsMap),
			},
			Expected: TSType("Record<string, Tag>"),
		},
		{
			Name: "map",
			Source: reflect.StructField{
				Name: "Milestones",
				Type: reflect.TypeOf(p.Milestones),
			},
			Expected: TSType("Record<number, number>"),
		},
		{
			Name: "sql.NullInt64",
			Source: reflect.StructField{
				Name: "Parent",
				Type: reflect.TypeOf(p.Parent),
			},
			Expected: TSType("number | null"),
		},
		{
			Name: "time",
			Source: reflect.StructField{
				Name: "CreatedAt",
				Type: reflect.TypeOf(p.CreatedAt),
			},
			Expected: TypeNumber,
		},
		{
			Name: "time pointer",
			Source: reflect.StructField{
				Name: "DeletedAt",
				Type: reflect.TypeOf(p.DeletedAt),
			},
			Expected: TSType("number | null"),
		},
		{
			Name: "single string type",
			Source: reflect.StructField{
				Name: "name",
				Type: reflect.TypeOf(name),
			},
			Expected: TypeString,
		},
	}

	tg := NewTypeGenerator(DefaultTypeGeneratorOpts)
	for _, tt := range tests {
		got := tg.getType(tt.Source)
		if got != tt.Expected {
			t.Errorf("`%s`: got %v, want %v", tt.Name, got, tt.Expected)
		}
	}
}

func Test_GetTypeWithObjectExpansion(t *testing.T) {
	p := Person{}

	tests := []struct {
		Name     string
		Source   reflect.StructField
		Expected TSType
	}{
		{
			Name: "struct (expanded)",
			Source: reflect.StructField{
				Name: "Address",
				Type: reflect.TypeOf(p.Address),
			},
			Expected: TSType(`{
        street?: string;
        city?: string;
        post_code: Capitalize<string>;
    }`),
		},
		{
			Name: "embedded struct (expanded)",
			Source: reflect.StructField{
				Name: "SocialMedia",
				Type: reflect.TypeOf(p.SocialMedia),
			},
			Expected: TSType(`{
        instagram_handle: string | null;
        twitter_handle: string;
        blogs: string[];
        hackernews_items: Array<number | null>;
    }`),
		},
	}

	tg := NewTypeGenerator(TypeGeneratorOpts{
		ExpandStruct: true,
	})

	for _, tt := range tests {
		got := tg.GetFieldType(tt.Source)
		if got != tt.Expected {
			t.Errorf("`%s`: got\n%v, want\n%v", tt.Name, got, tt.Expected)
		}
	}
}
