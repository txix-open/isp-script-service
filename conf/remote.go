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
	LogLevel                 log.Level          `schemaGen:"logLevel" schema:"Уровень логирования"`
	Scripts                  []ScriptDefinition `schema:"Скрипты"`
	CustomScripts            []ScriptDefinition `schema:"Кастомные скрипты"`
	SharedScript             string             `schema:"Общий скрипт"`
	ScriptExecutionTimeoutMs int                `validate:"required" schema:"Время выполнения скрипта"`
}

type ScriptDefinition struct {
	Id     string `validate:"required"`
	Script string `validate:"required"`
}
