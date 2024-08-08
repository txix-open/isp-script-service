package conf

import (
	"reflect"

	"github.com/txix-open/isp-kit/log"
	"github.com/txix-open/isp-kit/rc/schema"
	"github.com/txix-open/jsonschema"
)

// nolint:gochecknoinits
func init() {
	schema.CustomGenerators.Register("logLevel", func(field reflect.StructField, t *jsonschema.Schema) {
		t.Type = "string"
		t.Enum = []interface{}{"debug", "info", "error", "fatal"}
	})
}

type Remote struct {
	LogLevel                 log.Level `schemaGen:"logLevel" schema:"Уровень логирования"`
	Scripts                  []ScriptDefinition
	SharedScript             string
	ScriptExecutionTimeoutMs int `validate:"required"`
}

type ScriptDefinition struct {
	Id     string `validate:"required"`
	Script string `validate:"required"`
}
