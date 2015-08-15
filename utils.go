package v7

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func Int(reply interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case JSNumber:
		x := int(reply)
		if JSNumber(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	default:
		return 0, fmt.Errorf("go-v7: unexpected type for Int, got type %T", reply)
	}
}

func Int64(reply interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case JSNumber:
		x := int64(reply)
		if JSNumber(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	default:
		return 0, fmt.Errorf("go-v7: unexpected type for Int64, got type %T", reply)
	}
}

func Uint64(reply interface{}, err error) (uint64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		if reply < 0 {
			return 0, errNegativeInt
		}
		return uint64(reply), nil
	case JSNumber:
		if reply < 0.0 {
			return 0, errNegativeInt
		}
		return uint64(reply), nil
	default:
		return 0, fmt.Errorf("go-v7: unexpected type for Uint64, got type %T", reply)
	}
}

func Float64(reply interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case JSNumber:
		return float64(reply), nil
	default:
		return 0, fmt.Errorf("go-v7: unexpected type for Float64, got type %T", reply)
	}
}

func String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case JSString:
		return string(reply), nil
	case JSNumber:
		return strconv.FormatFloat(float64(reply), 'f', 6, 64), nil
	case bool:
		if reply {
			return "true", nil
		} else {
			return "false", nil
		}
	case []byte:
		return string(reply), nil
	default:
		return "", fmt.Errorf("go-v7: unexpected type for String, got type %T", reply)
	}
}

func Bool(reply interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	switch reply := reply.(type) {
	case bool:
		return bool(reply), nil
	default:
		return false, fmt.Errorf("go-v7: unexpected type for Bool, got type %T", reply)
	}
}

func Array(reply interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []byte:
		var arr []interface{}
		if err := json.Unmarshal(reply, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("go-v7: unexpected type for Array, got type %T", reply)
	}
}

func Object(reply interface{}, err error) (map[string]interface{}, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []byte:
		var obj map[string]interface{}
		if err := json.Unmarshal(reply, &obj); err != nil {
			return nil, err
		}
		return obj, nil
	default:
		return nil, fmt.Errorf("go-v7: unexpected type for Object, got type %T", reply)
	}
}
