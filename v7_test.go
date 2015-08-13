package v7_test

import (
	"github.com/bluele/go-v7"
	"log"
	"testing"
)

func newVM() *v7.V7 {
	return v7.New()
}

func TestExecBoolean(t *testing.T) {
	v := newVM()
	defer v.Destroy()
	ret, err := v.Exec(`true`)
	if err != nil {
		t.Error(err)
	}
	log.Println(ret)
	if bl, ok := ret.(bool); !ok || bl != true {
		t.Error("value: not true")
	}
	ret, err = v.Exec(`false`)
	if err != nil {
		t.Error(err)
	}
	if bl, ok := ret.(bool); !ok || bl != false {
		t.Error("value: not false")
	}
}
