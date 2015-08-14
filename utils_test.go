package v7_test

import (
	"github.com/bluele/go-v7"
	"testing"
)

func TestIntType(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	_, err := v7.Int(vm.Exec(`1`))
	if err != nil {
		t.Error("value should be 1.")
	}

	_, err = v7.Int(vm.Exec(`"golang"`))
	if err == nil {
		t.Error("value should be not int type.")
	}

	_, err = v7.Int(vm.Exec(`1.1`))
	if err == nil {
		t.Error("parser should be not int type.")
	}
}

func TestStringType(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	_, err := v7.String(vm.Exec(`"golang"`))
	if err != nil {
		t.Error("value should be string type.")
	}

	_, err = v7.String(vm.Exec(`1.1`))
	if err != nil {
		t.Error("value should be string type.")
	}

	_, err = v7.String(vm.Exec(`true`))
	if err != nil {
		t.Error("value should be string type.")
	}

	_, err = v7.String(vm.Exec(`[1,2,3]`))
	if err != nil {
		t.Error("value should be string type.")
	}
}

func TestBoolType(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	_, err := v7.Bool(vm.Exec(`true`))
	if err != nil {
		t.Error("value should be bool type.")
	}

	_, err = v7.Bool(vm.Exec(`false`))
	if err != nil {
		t.Error("value should be bool type.")
	}

	_, err = v7.Bool(vm.Exec(`"true"`))
	if err == nil {
		t.Error("value should be not bool type.")
	}
}

func TestArrayType(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	_, err := v7.Array(vm.Exec(`[1,"2",3]`))
	if err != nil {
		t.Error("value should be array type.")
	}

	_, err = v7.Array(vm.Exec(`10`))
	if err == nil {
		t.Error("value should be not array type.")
	}
}

func TestObjectType(t *testing.T) {
	vm := newVM()
	defer vm.Destroy()

	_, err := v7.Object(vm.Exec(`
  (function() {
    return {a: 10};
  })();`))
	if err != nil {
		t.Errorf("value should be object type. %v", err)
	}
}
