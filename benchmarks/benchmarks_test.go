package tests

import (
	"github.com/bluele/go-v7"
	"github.com/olebedev/go-duktape"
	"github.com/robertkrimen/otto"
	"testing"
)

const source = `
(function() {
  var a = 100;
  return a + 10;
})();
`

func BenchmarkV7(b *testing.B) {
	vm := v7.New()
	defer vm.Destroy()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := vm.Exec(source)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkDuktape(b *testing.B) {
	vm := duktape.NewContext()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if vm.PevalString(source) == 1 {
			panic(vm.SafeToString(-1))
		}
		vm.SafeToString(-1)
		vm.Pop()
	}
}

func BenchmarkOtto(b *testing.B) {
	vm := otto.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := vm.Run(source)
		if err != nil {
			panic(err)
		}
	}
}
