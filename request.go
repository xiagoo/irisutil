package irisutil

import (
	"reflect"
	"fmt"
	"gopkg.in/kataras/iris.v6"
	"strconv"
)

//FormRequest form
func  FormRequest(ctx *iris.Context, request interface{}, callback func() (interface{}, error)) (interface{}, error) {
	formValues := ctx.FormValues()
	if formValues == nil {
		panic(fmt.Sprintf("[FormRequest] fail to read post form"))
	}

	formMap := make(map[string]string)
	for k, v := range formValues {
		if len(v) > 0 {
			formMap[k] = v[0]
		} else {
			panic(fmt.Sprintf("[FormRequest] value is empty %#v", v))
		}
	}

	if len(formMap) == 0 {
		panic(fmt.Sprintf("[FormRequest] form is empty %#v", formMap))
	}

	reqType := reflect.TypeOf(request)
	if reqType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("[FormRequest] request must be  ptr type is %#v", reqType.Kind()))
	}

	fields := reqType.Elem()
	if fields.Kind() != reflect.Struct {
		panic(fmt.Sprintf("[FormRequest] request must be  struct type is %#v", fields.Kind()))
	}
	reqValue := reflect.ValueOf(request)

	for i := 0; i < fields.NumField(); i++ {
		f := fields.Field(i)

		fType := f.Type
		fTag := f.Tag.Get("json")

		if fTag == "" {
			panic(fmt.Sprintf("[FormRequest] struct need json tag"))
		}

		if formMap[fTag] == "" {
			fmt.Printf("[FormRequest] formMap[fTag] is empty , key = %v", fTag)
			continue
		}

		switch fType.Kind() {
		case reflect.String:
			reqValue.Elem().Field(i).SetString(formMap[fTag])

		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			formInt64, err := strconv.ParseInt(formMap[fTag], 10, 64)
			if err != nil {
				fmt.Printf("[FormRequest] formMap[fTag] = %#v convert to int64 failed", formMap[fTag])
				break
			}
			reqValue.Elem().Field(i).SetInt(formInt64)

		case reflect.Float32, reflect.Float64:
			formFloat64, err := strconv.ParseFloat(formMap[fTag], 64)
			if err != nil {
				fmt.Printf("[FormRequest] formMap[fTag] = %#v convert to float64 failed", formMap[fTag])
				break
			}
			reqValue.Elem().Field(i).SetFloat(formFloat64)

		case reflect.Bool:
			formBool, err := strconv.ParseBool(formMap[fTag])
			if err != nil {
				fmt.Printf("[FormRequest] formMap[fTag] = %#v convert to bool failed", formMap[fTag])
				break
			}
			reqValue.Elem().Field(i).SetBool(formBool)

		default:
			fmt.Printf("[FormRequest] formMap[fTag] is unknown type  %#v", formMap[fTag])
		}
	}
	return callback()
}
