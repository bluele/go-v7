# go-v7

go-v7 is golang binding to [v7](https://github.com/cesanta/v7).

# [v7](https://github.com/cesanta/v7)

V7 is a JavaScript engine written in C. It makes it possible to program Internet of Things (IoT) devices in JavaScript.

# Install

```
$ go get -u github.com/bluele/go-v7
```

# Example

## Run something in the VM

```go
package main

import (
  "fmt"
  "github.com/bluele/go-v7"
)

func main() {
  res, err = vm.Exec(`
    (function() {
      return 3.14;
    })();`)
  if err != nil {
    panic(err)
  }
  fmt.Println("result:", res)
}
```

Output:
```
result: 3.14
```

# Benchmarks

```
$ cd benchmarks && go test -bench .
testing: warning: no tests to run
PASS
BenchmarkV7       200000              9176 ns/op
BenchmarkDuktape           30000             56500 ns/op
BenchmarkOtto      50000             24019 ns/op
ok      github.com/bluele/go-v7/benchmarks      5.717s
```

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>