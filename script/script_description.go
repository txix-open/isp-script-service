package script

var scriptEngine = &Goja{}

func Default() Engine {
	return scriptEngine
}

type Engine interface {
	Compile(string) (Script, error)
	Execute(Script, interface{}) (interface{}, error)
}

type Script interface {
	Src() interface{}
}
