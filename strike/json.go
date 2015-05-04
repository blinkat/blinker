package strike

import (
	"fmt"
	"github.com/blinkat/blinker/strike/json"
	"reflect"
	"strings"
)

const (
	ConverJsonFormat_None       = iota // no format
	ConverJsonFormat_Upper             // all upper
	ConverJsonFormat_Lower             // all lower
	ConverJsonFormat_FirstUpper        // first word to upper
	ConverJsonFormat_FirstLower        // first word to upper
)

func ConverJsonFormat(obj interface{}, format int) string {
	kind := reflect.TypeOf(obj).Kind()

	if kind != reflect.Struct && kind != reflect.Map && kind != reflect.Ptr && kind != reflect.Interface {
		ret := conver_startement(obj, format)
		j := make(json.Json, 0)
		j = append(j, &json.JsonBlock{
			Name:  "Value",
			Value: ret,
		})
		return j.String()
	}
	return conver_startement(obj, format).String()
}

func ConverJson(obj interface{}) string {
	return ConverJsonFormat(obj, ConverJsonFormat_None)
}

func conver_startement(obj interface{}, format int) json.JsonValue {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Array, reflect.Slice:
		return conver_array(obj, format)
	case reflect.Map:
		return conver_map(obj, format)
	case reflect.Bool:
		return json.JsonBool(obj.(bool))

	case reflect.Float32:
		return json.JsonNumber(obj.(float32))
	case reflect.Float64:
		return json.JsonNumber(obj.(float64))
	case reflect.Int:
		return json.JsonNumber(obj.(int))
	case reflect.Int16:
		return json.JsonNumber(obj.(int16))
	case reflect.Int32:
		return json.JsonNumber(obj.(int32))
	case reflect.Int64:
		return json.JsonNumber(obj.(int64))
	case reflect.Int8:
		return json.JsonNumber(obj.(int8))

	case reflect.String:
		return json.JsonString(obj.(string))

	case reflect.Uint:
		return json.JsonNumber(obj.(int))
	case reflect.Uint16:
		return json.JsonNumber(obj.(uint16))
	case reflect.Uint32:
		return json.JsonNumber(obj.(uint32))
	case reflect.Uint64:
		return json.JsonNumber(obj.(uint64))
	case reflect.Uint8:
		return json.JsonNumber(obj.(uint8))

	case reflect.Struct:
		return conver_struct(obj, format)
	case reflect.Ptr, reflect.Interface:
		return conver_struct_ptr(obj, format)
	}
	fmt.Println("not type:", obj)
	return nil
}

func conver_array(obj interface{}, format int) json.JsonValue {
	ret := make(json.JsonArray, 0)
	value := reflect.ValueOf(obj)
	for i := 0; i < value.Len(); i++ {
		ret = append(ret, conver_startement(value.Index(i).Interface(), format))
	}
	return ret
}

func conver_map(obj interface{}, format int) json.JsonValue {
	ret := make(json.Json, 0)
	value := reflect.ValueOf(obj)
	keys := value.MapKeys()
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		item := value.MapIndex(key)

		ret = append(ret, &json.JsonBlock{
			Name:  handle_format(fmt.Sprint(key.Interface()), format),
			Value: conver_startement(item.Interface(), format),
		})
	}
	return ret
}

func conver_struct(obj interface{}, format int) json.JsonValue {
	return conver_struct_fields(reflect.ValueOf(obj), reflect.TypeOf(obj), format)
}

func conver_struct_ptr(obj interface{}, format int) json.JsonValue {
	return conver_struct_fields(reflect.ValueOf(obj).Elem(), reflect.TypeOf(obj).Elem(), format)
}

func conver_struct_fields(value reflect.Value, typeof reflect.Type, format int) json.JsonValue {
	ret := make(json.Json, 0)
	field_len := value.NumField()
	for i := 0; i < field_len; i++ {
		field := value.Field(i)
		if field.CanSet() {
			ret = append(ret, &json.JsonBlock{
				Name:  handle_format(typeof.Field(i).Name, format),
				Value: conver_startement(field.Interface(), format),
			})
		}
	}
	return ret
}

func handle_format(t string, f int) string {
	if t == "" {
		return t
	} else if f == ConverJsonFormat_Lower {
		return strings.ToLower(t)
	} else if f == ConverJsonFormat_Upper {
		return strings.ToUpper(t)
	} else if f == ConverJsonFormat_FirstLower {
		rs := []rune(t)
		return strings.ToUpper(string(rs[0])) + string(rs[1:])
	} else if f == ConverJsonFormat_FirstUpper {
		rs := []rune(t)
		return strings.ToLower(string(rs[0])) + string(rs[1:])
	} else {
		return t
	}
}
