package compile

import (
	"fmt"
	"github.com/integration-system/isp-lib/config"
	log "github.com/integration-system/isp-log"
	"isp-script-service/conf"
	"isp-script-service/log_code"
	"isp-script-service/script"
)

var (
	Script        compileScript
	remoteScripts = make(map[string]script.Script)
)

type compileScript struct{}

func (compileScript) GetById(id string) script.Script {
	return remoteScripts[id]
}

func (compileScript) Init(scriptDef []conf.ScriptDefinition) {
	for _, value := range scriptDef {
		prog, err := Script.Create(value.Script)
		if err != nil {
			log.Error(log_code.ErrorCreateScript, err)
			continue
		}
		remoteScripts[value.Id] = prog
	}
}

func (compileScript) Create(scr string) (script.Script, error) {
	cfg := config.GetRemote().(*conf.RemoteConfig)
	scr = fmt.Sprintf("%s%s%s%s", cfg.SharedScript, "(function() {", scr, "})();")
	return script.Default().Compile(scr)
}
