package conf

import "github.com/integration-system/isp-lib/structure"

type RemoteConfig struct {
	Scripts                  []ScriptDefinition
	SharedScript             string
	ScriptExecutionTimeoutMs int `valid:"required~Required"`
	Metrics                  structure.MetricConfiguration
}

type ScriptDefinition struct {
	Id     string `valid:"required~Required"`
	Script string `valid:"required~Required"`
}
