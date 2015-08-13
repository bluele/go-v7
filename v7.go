package v7

/*
#include "v7.h"
#include <stdlib.h>
int val_type(struct v7 *v7, uint64_t);
struct v7_function *v7_to_function(v7_val_t v);
char *v7_to_json(struct v7 *, v7_val_t val, char *buf, size_t buf_len);

int _val_type(struct v7 *v7, uint64_t t) {
  return val_type(v7, (v7_val_t)t);
}
const char *_v7_to_string(struct v7 *v, v7_val_t *value) {
  size_t string_len;
  return v7_to_string(v, value, &string_len);
}
*/
import "C"

import (
	"encoding/json"
	"errors"
	"unsafe"
)

type Value struct{}

var (
	UndefinedValue = &Value{}
	NullValue      = &Value{}
)

const (
	/* Primitive types */
	_V7_TYPE_UNDEFINED = iota
	_V7_TYPE_NULL
	_V7_TYPE_BOOLEAN
	_V7_TYPE_NUMBER
	_V7_TYPE_STRING
	_V7_TYPE_FOREIGN
	_V7_TYPE_CFUNCTION

	/* Different classes of Object type */
	_V7_TYPE_GENERIC_OBJECT
	_V7_TYPE_BOOLEAN_OBJECT
	_V7_TYPE_STRING_OBJECT
	_V7_TYPE_NUMBER_OBJECT
	_V7_TYPE_FUNCTION_OBJECT
	_V7_TYPE_CFUNCTION_OBJECT
	_V7_TYPE_REGEXP_OBJECT
	_V7_TYPE_ARRAY_OBJECT
	_V7_TYPE_DATE_OBJECT
	_V7_TYPE_ERROR_OBJECT
	_V7_TYPE_MAX_OBJECT_TYPE
	_V7_NUM_TYPES
)

// TODO fix memory leaks.
func _v7_to_json(ctx *C.struct_v7, val C.v7_val_t) []byte {
	size := 255
	buf := make([]byte, size)
	p := C.CString(string(buf))
	ret := C.v7_to_json(ctx, val, p, C.size_t(size))
	return []byte(C.GoString(ret))
}

type V7 struct {
	ctx *C.struct_v7
}

type Context struct {
	vm *C.struct_v7
}

type Function struct {
	ctx  *C.struct_v7
	repl C.v7_val_t
}

func (fc Function) Call(args ...interface{}) (interface{}, error) {
	result := C.v7_apply(fc.ctx, fc.repl, C.v7_create_undefined(), C.v7_create_undefined())
	return toValue(fc.ctx, result)
}

func New() *V7 {
	v := C.v7_create()
	return &V7{ctx: v}
}

func (v *V7) Exec(js string) (interface{}, error) {
	_js := C.CString(js)
	defer C.free(unsafe.Pointer(_js))

	var result C.v7_val_t
	C.v7_exec(v.ctx, &result, _js)

	return toValue(v.ctx, result)
}

func toValue(ctx *C.struct_v7, result C.v7_val_t) (interface{}, error) {
	switch C._val_type(ctx, C.uint64_t(result)) {
	case _V7_TYPE_UNDEFINED:
		return UndefinedValue, nil
	case _V7_TYPE_NULL:
		return NullValue, nil
	case _V7_TYPE_NUMBER:
		return float64(C.v7_to_number(result)), nil
	case _V7_TYPE_STRING:
		return C.GoString(C._v7_to_string(ctx, &result)), nil
	case _V7_TYPE_BOOLEAN:
		if int(C.v7_to_boolean(result)) == 0 {
			return false, nil
		} else {
			return true, nil
		}
	case _V7_TYPE_ARRAY_OBJECT:
		var a []interface{}
		js := _v7_to_json(ctx, result)
		if err := json.Unmarshal(js, &a); err != nil {
			return nil, err
		}
		return a, nil
	case _V7_TYPE_GENERIC_OBJECT:
		var m map[string]interface{}
		js := _v7_to_json(ctx, result)
		if err := json.Unmarshal(js, &m); err != nil {
			return nil, err
		}
		return m, nil
	case _V7_TYPE_FUNCTION_OBJECT:
		return Function{ctx, result}, nil
	default:
		return nil, errors.New("Undefined error")
	}
}

func (v *V7) Destroy() {
	C.v7_destroy(v.ctx)
}
