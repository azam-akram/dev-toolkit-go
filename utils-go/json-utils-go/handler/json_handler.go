package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/dev-toolkit-go/utils-go/json-utils-go/logger"
)

var handler Handler

type JsonHandler struct {
	logger logger.Logger
}

func NewJsonHandler() Handler {
	if handler == nil {
		handler = &JsonHandler{
			logger: logger.Init("info"),
		}
	}
	return handler
}

func (jh JsonHandler) GetLogger() logger.Logger {
	if jh.logger == nil {
		jh.logger = logger.Get()
	}
	return jh.logger
}

func (jh JsonHandler) ConvertGenericInterfaceToMap() {
	b := []byte(`{"k1":"v1","k2":6,"k3":["v3","v4"]}`)
	//fmt.Println(b)
	var i interface{}
	_ = json.Unmarshal(b, &i)
	fmt.Println(i)

	d := i.(map[string]interface{})

	for k, v := range d {
		switch vv := v.(type) {
		case string:
			fmt.Printf("key = %s, value = %s, value type = string\n", k, vv)
		case float64:
			fmt.Printf("key = %s, value = %f, value type = float64\n", k, vv)
		case []interface{}:
			fmt.Println(k, "'s value is a array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "unknown type")
		}
	}
}

func (jh JsonHandler) StringToStruct(s string, i interface{}) error {
	err := json.Unmarshal([]byte(s), i)
	if err != nil {
		return err
	}

	return nil
}

func (jh JsonHandler) StructToString(i interface{}) (string, error) {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (jh JsonHandler) StringToMap(s string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (jh JsonHandler) MapToString(m map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (jh JsonHandler) BytesToString(jsonBytes []byte) string {
	return string(jsonBytes)
}

func (jh JsonHandler) StringToBytes(s string) []byte {
	return []byte(s)
}

func (jh JsonHandler) StructToBytes(i interface{}) (jsonBytes []byte, err error) {
	jsonBytes, err = json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func (jh JsonHandler) BytesToStruct(b []byte, d interface{}) error {
	err := json.Unmarshal([]byte(b), &d)
	if err != nil {
		return err
	}

	return nil
}

func (jh JsonHandler) ModifyInputJson(s string) (map[string]interface{}, error) {
	var mapToProcess = make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &mapToProcess); err != nil {
		return nil, errors.New("cannot convert string to map")
	}

	jh.logger.Info("Before modification", "mapToProcess", mapToProcess)
	mapToProcess["degree"] = "phd"
	jh.logger.Info("After adding a new key-value", "mapToProcess", mapToProcess)

	return mapToProcess, nil
}

func (jh JsonHandler) DisplayAllJsonHandlers(str string) {
	//h.ConvertGenericInterfaceToMap()

	var emp Employee
	err := jh.StringToStruct(str, &emp)
	if err != nil {
		jh.GetLogger().Error("Error in converting string to struct", "error", err)
	}

	jh.GetLogger().Info("DisplayAllJsonHandlers::StringToStruct", "emp", emp)

	str, err = jh.StructToString(emp)
	if err != nil {
		jh.GetLogger().Info("Error in converting struct to string", "emp", emp)
	}
	jh.GetLogger().Info("DisplayAllJsonHandlers::ConvertStructToString", "str", str)

	jMap, err := jh.StringToMap(str)
	if err != nil {
		jh.GetLogger().Error("Error in StringToMap", "error", err)
	}
	jh.GetLogger().Info("DisplayAllJsonHandlers::ConvertStringToMap", "jMap", jMap)

	mapData := map[string]interface{}{
		"id":   "The ID",
		"user": "The User",
	}

	mapStr, err := jh.MapToString(mapData)
	if err != nil {
		jh.GetLogger().Error("Error in MapToString", "error", err)
	}

	jh.GetLogger().Info("DisplayAllJsonHandlers::ConvertMapToString", "mapStr", mapStr)

	jsonBytes := jh.StringToBytes(mapStr)
	jh.GetLogger().Info("DisplayAllJsonHandlers::ConvertStringToByte", "jsonBytes", jsonBytes)

	bytesStr := jh.BytesToString(jsonBytes)
	jh.GetLogger().Info("DisplayAllJsonHandlers::ConvertByteToString", "bytesStr", bytesStr)

	jsonBytes, err = jh.StructToBytes(emp)
	if err != nil {
		jh.GetLogger().Error("Error in StructToBytes", "error", err)
	}
	jh.GetLogger().Info("DisplayAllJsonHandlers::ConvertStructToByte", "jsonBytes", jsonBytes)

	err = jh.BytesToStruct(jsonBytes, &emp)
	if err != nil {
		jh.GetLogger().Error("Error in BytesToStruct", "error", err)
	}
	jh.GetLogger().Info("DisplayAllJsonHandlers::ConvertByteToStruct", "emp", emp)

	modifiedEmpMap, err := jh.ModifyInputJson(str)
	if err != nil {
		jh.GetLogger().Error("Error in ModifyInputJson", "error", err)
	}

	jh.GetLogger().Info("DisplayAllJsonHandlers::ModifyInputJson", "modifiedEmpMap", modifiedEmpMap)
}
