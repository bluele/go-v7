package v7_test

import (
	"github.com/bluele/go-v7"
	"testing"
)

func newVM() *v7.V7 {
	return v7.New()
}

func TestExecBoolean(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()
	ret, err := vm.Exec(`true`)
	if err != nil {
		t.Error(err)
	}
	if bl, ok := ret.(bool); !ok || bl != true {
		t.Error("value: not true")
	}
	ret, err = vm.Exec(`false`)
	if err != nil {
		t.Error(err)
	}
	if bl, ok := ret.(bool); !ok || bl != false {
		t.Error("value: not false")
	}
}

func TestJSFunctionString(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	res, err := vm.Exec(`(function() {
			return function(v) { return v + ":suffix"; };
		})();`)
	if err != nil {
		t.Error(err)
		return
	}

	arg := "test"
	val, err := v7.String(res.(*v7.JSFunction).Call(arg))
	if err != nil {
		t.Error(err)
	}
	if val != arg+":suffix" {
		t.Errorf("v should be %v", arg+":suffix")
	}
}

func TestJSFunctionInt64(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	res, err := vm.Exec(`(function() {
			return function(v) { return v * 3; };
		})();`)
	if err != nil {
		t.Error(err)
		return
	}

	val, err := v7.Int64(res.(*v7.JSFunction).Call(10))
	if err != nil {
		t.Error(err)
	}
	if val != 10*3 {
		t.Errorf("v should be %v", 10*3)
	}
}

func TestJSFunctionFloat64(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	res, err := vm.Exec(`(function() {
			return function(v) { return v * 3; };
		})();`)
	if err != nil {
		t.Error(err)
		return
	}

	val, err := v7.Float64(res.(*v7.JSFunction).Call(10.0))
	if err != nil {
		t.Error(err)
	}
	if val != 10.0*3 {
		t.Errorf("v should be %v", 10.0*3)
	}
}

func TestJSFunctionBool(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	res, err := vm.Exec(`(function() {
			return function(v) {
				if (v % 2 === 0) {
					return true;
				} else {
					return false;
				}
			};
		})();`)
	if err != nil {
		t.Error(err)
		return
	}

	val, err := v7.Bool(res.(*v7.JSFunction).Call(1))
	if err != nil {
		t.Error(err)
	}
	if val {
		t.Errorf("v should be false")
	}

	val, err = v7.Bool(res.(*v7.JSFunction).Call(2))
	if err != nil {
		t.Error(err)
	}
	if !val {
		t.Errorf("v should be true")
	}
}
