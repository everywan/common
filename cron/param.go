package cron

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// 按照 json 格式解析参数, 推荐.
func ParseJsonParams(paramsStr string, out interface{}) error {
	if strings.HasPrefix(paramsStr, "{") && strings.HasSuffix(paramsStr, "}") {
		if err := json.Unmarshal([]byte(paramsStr), out); err != nil {
			return errors.Wrapf(err, "unmarshal xxljob param error. param:%s", paramsStr)
		}
	}
	return nil
}

// 按照 urlquery 格式解析参数: param1=xxx&param2=xxxx, 取结构体中 tag:json 为列名.
func ParseUrlParams(paramStr string, param interface{}) {
	args := strings.Split(paramStr, "&")
	argMap := map[string]string{}
	for _, arg := range args {
		kvs := strings.SplitN(arg, "=", 2)
		if len(kvs) < 2 {
			continue
		}
		argMap[kvs[0]] = kvs[1]
	}

	paramV := reflect.ValueOf(param).Elem()
	paramT := paramV.Type()
	for i := 0; i < paramT.NumField(); i++ {
		field := paramT.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}
		if arg, ok := argMap[tag]; ok {
			switch field.Type.Kind() {
			case reflect.String:
				paramV.Field(i).SetString(arg)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				vInt, err := strconv.ParseInt(arg, 10, 64)
				if err != nil {
					continue
				}
				paramV.Field(i).SetInt(vInt)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vFloat, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					continue
				}
				paramV.Field(i).SetFloat(vFloat)
			case reflect.Slice: // 只支持 []string
				value := strings.Split(arg, ",")
				paramV.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
}
