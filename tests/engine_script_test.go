// nolint:gochecknoglobals
package tests_test

import (
	"strings"
	"testing"

	json2 "github.com/txix-open/isp-kit/json"

	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	json "layeh.com/gopher-json"
)

const (
	jsScript = `
		var a = {};
		for (i=0; i<b.length; i++) {
			a[b[i]["id_tp_cd"]] = b[i]
		}
	`
	luaScript = `
		dst = {}
		for k,v in pairs(src) do
			dst[v.id_tp_cd] = v 
		end
`
)

var (
	recordsExample = []map[string]interface{}{
		{
			"ref_num":               "IIИК778068",
			"id_tp_cd":              "1000017",
			"start_dt":              "2006-10-10 00:00:00.000",
			"etalon_id":             "7d53bf34-0aa6-11e9-9799-87ba01ce01e5",
			"ext_upd_date":          "2018-01-11 17:52:01.000",
			"last_update_dt":        "2018-09-13 23:41:35.324",
			"identification_issuer": "Красногорским Управлением ЗАГС Главного Управления ЗАГС Московской области",
		},
		{
			"ref_num":        "18605520974",
			"id_tp_cd":       "1000015",
			"start_dt":       "2018-09-13 00:00:00.000",
			"etalon_id":      "7d57dde9-0aa6-11e9-9799-87ba01ce01e5",
			"ext_upd_date":   "2018-01-11 17:52:01.000",
			"last_update_dt": "2018-09-13 23:41:35.330",
		},
		{
			"ref_num":        "5090399778001137",
			"id_tp_cd":       "1000014",
			"start_dt":       "2018-09-13 23:41:35.327",
			"etalon_id":      "7d58a13e-0aa6-11e9-9799-87ba01ce01e5",
			"ext_upd_date":   "2018-01-11 17:52:01.000",
			"last_update_dt": "2018-09-13 23:41:35.327",
		},
	}
	srcString []byte
	expecting = make(map[string]interface{})
)

// nolint:gochecknoinits
func init() {
	for _, v := range recordsExample {
		expecting[v["id_tp_cd"].(string)] = v // nolint:forcetypeassert
	}
	if bytes, err := json2.Marshal(recordsExample); err != nil {
		panic(err)
	} else {
		srcString = bytes
	}
}

func TestOtto(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	prog, err := parser.ParseFile(nil, "example", jsScript, 0)
	require.NoError(err)

	runTime := otto.New()
	err = runTime.Set("b", recordsExample)
	require.NoError(err)

	_, err = runTime.Run(prog)
	require.NoError(err)

	res, err := runTime.Get("a")
	require.NoError(err)

	_, err = res.Export()
	require.NoError(err)
}

func TestGoja(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	prog := goja.MustCompile("test.js", jsScript, false)

	runTime := goja.New()
	err := runTime.Set("b", recordsExample)
	require.NoError(err)

	_, err = runTime.RunProgram(prog)
	require.NoError(err)

	a := runTime.Get("a")
	_ = a.Export()
}

func TestGopherLua(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	l := lua.NewState()
	defer l.Close()

	b, err := json2.Marshal(recordsExample)
	require.NoError(err)

	val, err := json.Decode(l, b)
	require.NoError(err)

	l.SetGlobal("src", val)

	err = l.DoString(luaScript)
	require.NoError(err)

	dst := l.GetGlobal("dst")
	b, err = json.Encode(dst)
	require.NoError(err)

	m := make(map[string]interface{})
	err = json2.Unmarshal(b, &m)
	require.NoError(err)

	require.EqualValues(expecting, m)
}

func BenchmarkOtto(b *testing.B) {
	prog, err := parser.ParseFile(nil, "example", jsScript, 0)
	if err != nil {
		return
	}

	for range b.N {
		runTime := otto.New()
		err = runTime.Set("b", recordsExample)
		if err != nil {
			return
		}

		_, err = runTime.Run(prog)
		if err != nil {
			return
		}

		res, err := runTime.Get("a")
		if err != nil {
			return
		}
		_, err = res.Export()
		if err != nil {
			return
		}
	}
}

func BenchmarkGoja(b *testing.B) {
	prog := goja.MustCompile("test.js", jsScript, false)

	for range b.N {
		runTime := goja.New()
		err := runTime.Set("b", recordsExample)
		if err != nil {
			return
		}

		_, err = runTime.RunProgram(prog)
		if err != nil {
			return
		}

		a := runTime.Get("a")
		_ = a.Export()
	}
}

func BenchmarkGopherLua(b *testing.B) {
	chunk, err := parse.Parse(strings.NewReader(luaScript), "test")
	if err != nil {
		panic(err)
	}
	proto, err := lua.Compile(chunk, "test")
	if err != nil {
		panic(err)
	}

	for range b.N {
		l := lua.NewState()

		val, err := json.Decode(l, srcString)
		if err != nil {
			panic(err)
		}
		l.SetGlobal("src", val)

		prog := l.NewFunctionFromProto(proto)
		l.Push(prog)
		err = l.PCall(0, lua.MultRet, nil)
		if err != nil {
			panic(err)
		}

		dst := l.GetGlobal("dst")
		c, err := json.Encode(dst)
		if err != nil {
			panic(err)
		}
		m := make(map[string]interface{})
		err = json2.Unmarshal(c, &m)
		if err != nil {
			panic(err)
		}
		l.Close()
	}
}
