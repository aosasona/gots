package generator

import "testing"

type Addr struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
}

type Student struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Age       int            `json:"age"`
	Email     string         `json:"email"`
	Contact   []string       `json:"contact"`
	Grades    map[string]int `json:"grades"`
	Addr      `json:"address"`
}

func Test_ObjectTypeGeneration(t *testing.T) {
	tests := []struct {
		Name     string
		Source   any
		Expected string
	}{
		{
			Name:   "student",
			Source: Student{},
			Expected: `export type Student = {
    first_name: string;
    last_name: string;
    age: number;
    email: string;
    contact: string[];
    grades: Record<string, number>;
    address: Addr;
}`,
		},
		{
			Name:   "address",
			Source: Addr{},
			Expected: `export type Addr = {
    street: string;
    city: string;
    postcode: string;
}`,
		},
	}

	g := NewGenerator(Opts{
		UseTypeForObjects: true,
	})

	for _, tt := range tests {
		got := g.Generate(tt.Source)
		if got != tt.Expected {
			t.Errorf("`%s`: got\n%v, want\n%v", tt.Name, got, tt.Expected)
		}
	}
}
