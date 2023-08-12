package config

import "fmt"

type Case string

const (
	CaseCamel  Case = "camel"
	CasePascal Case = "pascal"
	CaseSnake  Case = "snake"
)

type Config struct {
	Enabled           *bool
	OutputFile        *string
	UseTypeForObjects *bool
	Case              Case // to be implemented later
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

	if other.Case != "" {
		caseOrDefault = other.Case
	}

	return Config{
		Enabled:           &enabled,
		OutputFile:        &outputFile,
		UseTypeForObjects: &useTypeForObjects,
		Case:              caseOrDefault,
	}
}
