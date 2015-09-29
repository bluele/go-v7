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
	"fmt"
	"unsafe"
)

var (
	UndefinedValue = &JSUndefined{}
	NullValue      = &JSNull{}
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

	defaultBufSize = 255
)

// enum v7_err
const (
	_V7_OK = iota
	_V7_SYNTAX_ERROR
	_V7_EXEC_EXCEPTION
	_V7_STACK_OVERFLOW
	_V7_AST_TOO_LARGE
	_V7_INVALID_ARG
)

var (
	GO_V7_ERROR_TEMPLATE = "v7 error: %v"
	GO_V7_SYNTAX_ERROR   = fmt.Errorf(GO_V7_ERROR_TEMPLATE, _V7_SYNTAX_ERROR)
	GO_V7_EXEC_EXCEPTION = fmt.Errorf(GO_V7_ERROR_TEMPLATE, _V7_EXEC_EXCEPTION)
	GO_V7_STACK_OVERFLOW = fmt.Errorf(GO_V7_ERROR_TEMPLATE, _V7_STACK_OVERFLOW)
	GO_V7_AST_TOO_LARGE  = fmt.Errorf(GO_V7_ERROR_TEMPLATE, _V7_AST_TOO_LARGE)
	GO_V7_INVALID_ARG    = fmt.Errorf(GO_V7_ERROR_TEMPLATE, _V7_INVALID_ARG)
)

func _v7_to_json(ctx *C.struct_v7, val C.v7_val_t, size int) []byte {
	buf := make([]byte, size)
	p := C.CString(string(buf))
	defer C.free(unsafe.Pointer(p))
	ret := C.v7_to_json(ctx, val, p, C.size_t(size))
	return []byte(C.GoString(ret))
}

type V7 struct {
	ctx  *Context
	size int
}

type JSFunction struct {
	ctx   *Context
	value C.v7_val_t
}

type Context struct {
	rctx    *C.struct_v7
	bufSize int
}

func NewContext() *Context {
	return &Context{
		rctx:    C.v7_create(),
		bufSize: defaultBufSize,
	}
}

// Call JS function with arguments.
// Arguments is converted as following table.
// Go: Javascript
// uint, int, float: number
// []byte, string: string
// bool: boolean
func (fc *JSFunction) Call(args ...interface{}) (interface{}, error) {
	var cargs *C.v7_val_t

	if len(args) == 0 {
		_cargs := C.v7_create_undefined()
		cargs = &_cargs
	} else {
		var cobj *C.v7_val_t
		_cargs := C.v7_create_array(fc.ctx.rctx)
		for _, arg := range args {
			switch arg.(type) {
			case uint:
				obj := C.v7_create_number(C.double(arg.(uint)))
				cobj = &obj
			case uint32:
				obj := C.v7_create_number(C.double(arg.(uint32)))
				cobj = &obj
			case uint64:
				obj := C.v7_create_number(C.double(arg.(uint64)))
				cobj = &obj
			case int:
				obj := C.v7_create_number(C.double(arg.(int)))
				cobj = &obj
			case int32:
				obj := C.v7_create_number(C.double(arg.(int32)))
				cobj = &obj
			case int64:
				obj := C.v7_create_number(C.double(arg.(int64)))
				cobj = &obj
			case float32:
				obj := C.v7_create_number(C.double(arg.(float32)))
				cobj = &obj
			case float64:
				obj := C.v7_create_number(C.double(arg.(float64)))
				cobj = &obj
			case []byte, string:
				ptr := C.CString(arg.(string))
				defer C.free(unsafe.Pointer(ptr))
				obj := C.v7_create_string(fc.ctx.rctx, ptr, C.size_t(len(arg.(string))), 0)
				cobj = &obj
			case bool:
				obj := C.v7_create_boolean(arg.(C.int))
				cobj = &obj
			default:
				return nil, fmt.Errorf("Unsupported type %T", arg)
			}
			C.v7_array_push(fc.ctx.rctx, _cargs, *cobj)
		}
		cargs = &_cargs
	}

	var result C.v7_val_t
	err := wrapCError(C.v7_apply(fc.ctx.rctx, &result, fc.value, C.v7_create_undefined(), *cargs))
	if err != nil {
		return nil, err
	}

	return toValue(fc.ctx, result)
}

func wrapCError(cerr uint32) error {
	switch cerr {
	case _V7_OK:
		return nil
	case _V7_SYNTAX_ERROR:
		return GO_V7_SYNTAX_ERROR
	case _V7_EXEC_EXCEPTION:
		return GO_V7_EXEC_EXCEPTION
	case _V7_STACK_OVERFLOW:
		return GO_V7_STACK_OVERFLOW
	case _V7_AST_TOO_LARGE:
		return GO_V7_AST_TOO_LARGE
	case _V7_INVALID_ARG:
		return GO_V7_INVALID_ARG
	default:
		return fmt.Errorf(GO_V7_ERROR_TEMPLATE, cerr)
	}
}

func New() *V7 {
	return &V7{ctx: NewContext()}
}

func (v *V7) RawContext() *C.struct_v7 {
	return v.ctx.rctx
}

func (v *V7) BufferSize() int {
	return v.size
}

func (v *V7) ChangeBufferSize(size int) {
	v.ctx.bufSize = size
}

func (v *V7) Exec(js string) (interface{}, error) {
	_js := C.CString(js)
	defer C.free(unsafe.Pointer(_js))

	var result C.v7_val_t
	C.v7_exec(v.ctx.rctx, _js, &result)

	return toValue(v.ctx, result)
}

/*
Value compatibility between js and go.

js => go
number => float64
bool => boolean
string => string
function => v7.Function
array => []byte
object => []byte
*/
func toValue(ctx *Context, result C.v7_val_t) (interface{}, error) {
	rctx := ctx.rctx
	switch C._val_type(rctx, C.uint64_t(result)) {
	case _V7_TYPE_UNDEFINED:
		return UndefinedValue, nil
	case _V7_TYPE_NULL:
		return NullValue, nil
	case _V7_TYPE_NUMBER:
		return JSNumber(C.v7_to_number(result)), nil
	case _V7_TYPE_STRING:
		return JSString(C.GoString(C._v7_to_string(rctx, &result))), nil
	case _V7_TYPE_BOOLEAN:
		if int(C.v7_to_boolean(result)) == 0 {
			return JSFalse, nil
		} else {
			return JSTrue, nil
		}
	case _V7_TYPE_ARRAY_OBJECT, _V7_TYPE_GENERIC_OBJECT:
		return _v7_to_json(rctx, result, ctx.bufSize), nil
	case _V7_TYPE_FUNCTION_OBJECT:
		return &JSFunction{ctx, result}, nil
	default:
		return nil, errUndefinedType
	}
}

func (v *V7) Destroy() {
	C.v7_destroy(v.RawContext())
}
