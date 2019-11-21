package adapter

import (
	"encoding/json"
	"github.com/integration-system/isp-lib/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"isp-script-service/compile"
	"isp-script-service/conf"
	"isp-script-service/controller"
	"isp-script-service/domain"
	"isp-script-service/script"
	"log"
	"os"
	"testing"
)

const (
	remoteConf = `{
		"scriptExecutionTimeoutMs": 200
	}`

	id      = "123"
	version = 10

	objectFile = "obj.json"
	scriptFile = "custom_date.js"
)

var data map[string]interface{}
var request = domain.ExecuteRequest{}

func init() {
	config.InitRemoteConfig(&conf.RemoteConfig{}, []byte(remoteConf))

	obj, err := readFile(objectFile)
	if err != nil {
		return
	}
	if err = json.Unmarshal(obj, &data); err != nil {
		return
	}

	obj, err = readFile(scriptFile)
	if err != nil {
		return
	}
	script := string(obj)
	arg := map[string]interface{}{
		"id":      id,
		"version": version,
		"data":    data,
	}
	request = domain.ExecuteRequest{
		Script: script,
		Arg:    arg,
	}
}

func TestCustomDate(t *testing.T) {
	a := assert.New(t)

	execute := controller.Script.Execute(request)
	if !a.Zero(execute.Error) {
		return
	}
	executeResult, _ := json.Marshal(execute.Result)

	customDate := getCustomData(data, id, version)
	customDateResult, _ := json.Marshal(customDate)

	if !a.Equal(string(customDateResult), string(executeResult)) {
		log.Println("not equal")
	}
}

func BenchmarkGoCustomData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getCustomData(data, id, version)
		//json.Marshal(customDate)

	}
}

func BenchmarkJsCustomData(b *testing.B) {
	scr, err := compile.Script.Create(request.Script)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = script.Default().Execute(scr, request.Arg)
	}
}

func BenchmarkInlineJson(b *testing.B) {
	c, err := compile.Script.Create("var b = JSON.parse(arg); return JSON.stringify(b)")
	if err != nil {
		panic(err)
	}
	data, err := readFile(objectFile)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		_, err := script.Default().Execute(c, string(data))
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsoniterJson(b *testing.B) {
	c, err := compile.Script.Create("var b = arg; return b")
	if err != nil {
		panic(err)
	}
	data, err := readFile(objectFile)
	if err != nil {
		panic(err)
	}
	json := jsoniter.ConfigFastest
	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		if err := json.Unmarshal(data, &m); err != nil {
			panic(err)
		}
		res, err := script.Default().Execute(c, m)
		if err != nil {
			panic(err)
		}
		if _, err = json.Marshal(res); err != nil {
			panic(err)
		}
	}

}

func readFile(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	obj, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
