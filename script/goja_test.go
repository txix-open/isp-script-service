package script_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestGoja_AddFunction(t *testing.T) {
	a := assert.New(t)

	const SCRIPT = `
	var response = {};
	try {
		response["test"] = f("string", "unknown", "test");
	} catch (e) {
		if (!(e instanceof GoError)) {
			throw(e);
		}
		if (e.value.Error() !== "TEST") {
			throw("Unexpected value: " + e.value.Error());
		}
	}
	return response;
	`

	f := func(varchar string, integer int, object string) (interface{}, error) {
		a.Equal("string", varchar)
		a.Equal(0, integer)
		a.Equal("test", object)

		return "test", nil
	}
	vm := goja.New()
	vm.Set("f", f)
	resp, err := vm.RunString(fmt.Sprintf("(function() { %s })();", SCRIPT))
	a.NoError(err)
	a.Equal(resp.Export(), map[string]interface{}{"test": "test"})

	f2 := func(varchar string, integer int, object string) (interface{}, error) {
		a.Equal("string", varchar)
		a.Equal(0, integer)
		a.Equal("test", object)
		//nolint:goerr113
		return "test", errors.New("TEST")
	}
	vm = goja.New()
	vm.Set("f", f2)
	resp, err = vm.RunString(fmt.Sprintf("(function() { %s })();", SCRIPT))
	a.NoError(err)
	a.NoError(err)
	a.Equal(resp.Export(), map[string]interface{}{})
}
