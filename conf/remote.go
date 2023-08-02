package conf

import (
	"reflect"

	"github.com/integration-system/isp-kit/log"
	"github.com/integration-system/isp-kit/rc/schema"
	"github.com/integration-system/jsonschema"
)

// nolint:gochecknoinits
func init() {
	schema.CustomGenerators.Register("logLevel", func(field reflect.StructField, t *jsonschema.Type) {
		t.Type = "string"
		t.Enum = []interface{}{"debug", "info", "error", "fatal"}
	})
}

type Remote struct {
	LogLevel                 log.Level `schemaGen:"logLevel" schema:"Уровень логирования"`
	Scripts                  []ScriptDefinition
	SharedScript             string
	ScriptExecutionTimeoutMs int `valid:"required~Required"`
}

type ScriptDefinition struct {
	Id     string `valid:"required~Required"`
	Script string `valid:"required~Required"`
}
