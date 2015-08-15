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

func (fc *JSFunction) Call(args ...interface{}) (interface{}, error) {
	result := C.v7_apply(fc.ctx.rctx, fc.value, C.v7_create_undefined(), C.v7_create_undefined())
	return toValue(fc.ctx, result)
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
	C.v7_exec(v.ctx.rctx, &result, _js)

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
