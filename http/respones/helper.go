package respones

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Success(c echo.Context, params ...interface{}) error {
	res := make(map[string]interface{})
	res["code"] = 0
	res["msg"] = ""
	res["meta"] = nil
	res["data"] = []interface{}{}
	switch len(params) {
	case 0:
	case 1:
		res["msg"] = params[0].(string)
	case 2:
		res["msg"] = params[0].(string)
		res["meta"] = InterfaceToMap(params[1])
	case 3:
		res["msg"] = params[0]
		res["meta"] = params[1]
		res["data"] = params[2]
	}
	return c.JSON(200, res)

}
func Error(c echo.Context, code int, params ...interface{}) error {
	res := make(map[string]interface{})
	res["code"] = code
	res["msg"] = ""
	res["meta"] = nil
	res["data"] = []interface{}{}
	switch len(params) {
	case 0:
	case 1:
		res["msg"] = params[0].(string)
	case 2:
		res["msg"] = params[0].(string)
		res["meta"] = params[1].(map[string]interface{})
	case 3:
		res["msg"] = params[0]
		res["meta"] = params[1]
		res["data"] = params[2]
	}
	return c.JSON(200, res)

}
func InterfaceToMap(obj interface{}) map[string]interface{} {
	if meta, ok := obj.(map[string]interface{}); ok {
		return meta
	} else {
		marshal, _ := json.Marshal(obj)
		json.Unmarshal(marshal, &meta)
		return meta
	}
}
