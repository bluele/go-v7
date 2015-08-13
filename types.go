package v7

type (
	JSString string
	JSNumber float64
	JSObject map[string]interface{}
	JSArray  []interface{}

	JSUndefined struct{}
	JSNull      struct{}
)

const (
	JSTrue  = true
	JSFalse = false
)
