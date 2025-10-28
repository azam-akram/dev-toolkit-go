package handler

import (
	"encoding/json"
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

func (jh JsonHandler) StringToStruct(s string, obj any) error {
	err := json.Unmarshal([]byte(s), &obj)
	if err != nil {
		jh.logger.Error("StringToStruct: Failed to unmarshal json")
		return err
	}
	return nil
}

func (jh JsonHandler) StructToString(obj any) (string, error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		jh.logger.Error("StructToString: Failed to marshal struct", "error", err)
		return "", err
	}

	return string(objBytes), nil
}

func (jh JsonHandler) StringToMap(s string) (map[string]any, error) {
	var m map[string]any
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		jh.logger.Error("StringToMap: Failed to unmarshal json")
		return nil, err
	}

	return m, nil
}

func (jh JsonHandler) MapToString(m map[string]any) (string, error) {
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

func (jh JsonHandler) StructToBytes(i any) ([]byte, error) {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		jh.logger.Error("StringToMap: Failed to marshal object")
		return nil, err
	}

	return jsonBytes, nil
}

func (jh JsonHandler) BytesToStruct(b []byte, d any) error {
	err := json.Unmarshal([]byte(b), &d)
	if err != nil {
		return err
	}

	return nil
}

func (jh JsonHandler) ModifyInputJson(s string) (map[string]any, error) {
	var mapToProcess = make(map[string]any)
	if err := json.Unmarshal([]byte(s), &mapToProcess); err != nil {
		jh.logger.Error("ModifyInputJson: Failed to unmarshal object")
		return nil, err
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
		jh.logger.Error("Error in converting string to struct", "error", err)
	}

	jh.logger.Info("DisplayAllJsonHandlers::StringToStruct", "emp", emp)

	str, err = jh.StructToString(emp)
	if err != nil {
		jh.logger.Info("Error in converting struct to string", "emp", emp)
	}
	jh.logger.Info("DisplayAllJsonHandlers::ConvertStructToString", "str", str)

	jMap, err := jh.StringToMap(str)
	if err != nil {
		jh.logger.Error("Error in StringToMap", "error", err)
	}
	jh.logger.Info("DisplayAllJsonHandlers::ConvertStringToMap", "jMap", jMap)

	mapData := map[string]any{
		"id":   "The ID",
		"user": "The User",
	}

	mapStr, err := jh.MapToString(mapData)
	if err != nil {
		jh.logger.Error("Error in MapToString", "error", err)
	}

	jh.logger.Info("DisplayAllJsonHandlers::ConvertMapToString", "mapStr", mapStr)

	jsonBytes := jh.StringToBytes(mapStr)
	jh.logger.Info("DisplayAllJsonHandlers::ConvertStringToByte", "jsonBytes", jsonBytes)

	bytesStr := jh.BytesToString(jsonBytes)
	jh.logger.Info("DisplayAllJsonHandlers::ConvertByteToString", "bytesStr", bytesStr)

	jsonBytes, err = jh.StructToBytes(emp)
	if err != nil {
		jh.logger.Error("Error in StructToBytes", "error", err)
	}
	jh.logger.Info("DisplayAllJsonHandlers::ConvertStructToByte", "jsonBytes", jsonBytes)

	err = jh.BytesToStruct(jsonBytes, &emp)
	if err != nil {
		jh.logger.Error("Error in BytesToStruct", "error", err)
	}
	jh.logger.Info("DisplayAllJsonHandlers::ConvertByteToStruct", "emp", emp)

	modifiedEmpMap, err := jh.ModifyInputJson(str)
	if err != nil {
		jh.logger.Error("Error in ModifyInputJson", "error", err)
	}

	jh.logger.Info("DisplayAllJsonHandlers::ModifyInputJson", "modifiedEmpMap", modifiedEmpMap)
}

func (jh JsonHandler) ProcessGenericMap(d map[string]any) {
	for k, v := range d {
		switch vv := v.(type) {
		case string:
			jh.logger.Info("Data type found", "key", k, "type", "string", "value", vv)
		case float64:
			jh.logger.Info("Data type found", "key", k, "type", "float64", "value", vv)
		case bool:
			jh.logger.Info("Data type found", "key", k, "type", "bool", "value", vv)
		case []any:
			jh.logger.Info("Data type found", "key", k, "type", "array", "length", len(vv))
			for i, u := range vv {
				jh.logger.Info("Array element", "key", k, "index", i, "element_type", fmt.Sprintf("%T", u), "value", u)
			}
		case map[string]any:
			jh.logger.Info("Data type found", "key", k, "type", "nested_object")
		default:
			jh.logger.Error("Unknown type encountered", "key", k, "actual_type", fmt.Sprintf("%T", v))
		}
	}
}
