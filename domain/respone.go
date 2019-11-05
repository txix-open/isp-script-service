package domain

const (
	ErrorCompile = "Compilation"
	ErrorRunTime = "Runtime"
)

type ScriptResp struct {
	Result interface{}
	Error  *Error
}

type Error struct {
	Type        string
	Description string
}
