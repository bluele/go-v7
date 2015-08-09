package main

import (
	"fmt"
	"github.com/bluele/go-v7"
)

func main() {
	vm := v7.New()
	defer vm.Destroy()

	res, err := vm.Exec("100")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	res, err = v7.Bytes(vm.Exec(`"python"`))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	res, err = vm.Exec(`(function() {
			return function() { return 12; };
		})();`)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	fmt.Println(res.(v7.Function).Call())

	res, err = vm.Exec(`[1,2,3]`)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	res, err = v7.Float(vm.Exec(`
		(function() {
			return 100.1;
		})();`))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
