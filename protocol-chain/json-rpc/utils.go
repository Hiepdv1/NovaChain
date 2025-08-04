package jsonrpc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	log "github.com/sirupsen/logrus"
)

func CheckError(message string, err error) {
	if err != nil {
		log.Panic(fmt.Sprintf("%s:%v", message, err.Error()))
	}
}

func encodeBytes(data []byte) string {
	return fmt.Sprintf("%x", data)
}

func SafeMarshalJSON(v any) ([]byte, error) {
	return json.Marshal(processValue(v))
}

func processValue(v any) any {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}

		if b, ok := rv.Interface().(*big.Int); ok {
			return b.String()
		}

		return processValue(rv.Elem().Interface())

	case reflect.Struct:
		if b, ok := rv.Interface().(big.Int); ok {
			return b.String()
		}

		result := make(map[string]any)
		rt := rv.Type()
		for i := range rv.NumField() {
			field := rt.Field(i)
			name := field.Name
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				name = jsonTag
			}
			val := rv.Field(i).Interface()
			result[name] = processValue(val)
		}
		return result

	case reflect.Slice, reflect.Array:
		elemType := rv.Type().Elem().Kind()
		if elemType == reflect.Uint8 {
			return encodeBytes(rv.Bytes())
		}

		var arr []any
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			arr = append(arr, processValue(elem))
		}
		return arr

	case reflect.Map:
		result := make(map[string]any)
		for _, key := range rv.MapKeys() {
			val := rv.MapIndex(key).Interface()
			result[fmt.Sprintf("%v", key.Interface())] = processValue(val)
		}
		return result

	default:
		return v
	}
}
