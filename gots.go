package gots

import "reflect"

type MappedTSType string

const (
	_STRING  MappedTSType = "string"
	_BOOLEAN              = "boolean"

	_INVALID = ""
)

type gots struct {
	config Config
}

type Config struct {
	OutputDir             string
	OutputInMultipleFiles string
}

func New(config Config) *gots {
	return &gots{
		config,
	}
}

func getMappedType(t reflect.Kind) MappedTSType {
	switch t {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return _STRING
	default:
		return _INVALID
	}
}
