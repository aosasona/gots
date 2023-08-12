package config

import "fmt"

// Debug is a global variable that can be used to enable or disable debug mode
var Debug = false

func SetDebug(v bool) {
	Debug = v
}

type Case string

const (
	CaseCamel  Case = "camel"
	CasePascal Case = "pascal"
	CaseSnake  Case = "snake"
)

type Config struct {
	Enabled           *bool   // can be used to disable or enable the generation of types, defaults to false
	OutputFile        *string // if nil, will default to types.ts in the current directory
	UseTypeForObjects *bool   // if true, will use `type Foo = ...` instead of `interface Foo {...}`
	ExpandObjectTypes *bool   // if true, will expand object types instead of just using the name (e.g foo: { bar: string } instead of foo: Bar)
	Case              Case    // to be implemented later
}

func (c Config) EnabledOrDefault() bool {
	if c.Enabled == nil {
		fmt.Println("c.Enabled is nil, check your config")
		return false
	}

	return *c.Enabled
}

func (c Config) OutputFileOrDefault() string {
	if c.OutputFile == nil {
		return "types.ts"
	}

	return *c.OutputFile
}

func (c Config) ExpandObjectTypesOrDefault() bool {
	if c.ExpandObjectTypes == nil {
		return false
	}

	return *c.ExpandObjectTypes
}

func (c Config) UseTypeForObjectsOrDefault() bool {
	if c.UseTypeForObjects == nil {
		return false
	}

	return *c.UseTypeForObjects
}

func (c Config) CaseOrDefault() Case {
	if c.Case == "" {
		return CaseCamel
	}

	return c.Case
}

func (c Config) Merge(other Config) Config {
	enabled := c.EnabledOrDefault()
	outputFile := c.OutputFileOrDefault()
	useTypeForObjects := c.UseTypeForObjectsOrDefault()
	expandObjectTypes := c.ExpandObjectTypesOrDefault()
	caseOrDefault := c.CaseOrDefault()

	if other.Enabled != nil {
		enabled = *other.Enabled
	}

	if other.OutputFile != nil {
		outputFile = *other.OutputFile
	}

	if other.UseTypeForObjects != nil {
		useTypeForObjects = *other.UseTypeForObjects
	}

	if other.ExpandObjectTypes != nil {
		expandObjectTypes = *other.ExpandObjectTypes
	}

	if other.Case != "" {
		caseOrDefault = other.Case
	}

	return Config{
		Enabled:           &enabled,
		OutputFile:        &outputFile,
		UseTypeForObjects: &useTypeForObjects,
		ExpandObjectTypes: &expandObjectTypes,
		Case:              caseOrDefault,
	}
}
