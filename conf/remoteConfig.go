package conf

type RemoteConfig struct {
	Scripts                  []ScriptDefinition
	SharedScript             string
	ScriptExecutionTimeoutMs int `valid:"required~Required"`
}

type ScriptDefinition struct {
	Id     string `valid:"required~Required"`
	Script string `valid:"required~Required"`
}
