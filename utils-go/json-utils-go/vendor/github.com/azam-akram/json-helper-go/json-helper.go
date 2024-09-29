package jsonhelper

import (
	"encoding/json"
	"errors"
)

var handler Helper

type JsonHandler struct {
	logger Logger
}

func NewJsonHelper() Helper {
	if handler == nil {
		handler = &JsonHandler{
			logger: GetLogger(),
		}
	}
	return handler
}

func (jh JsonHandler) GetLogger() Logger {
	if jh.logger == nil {
		jh.logger = GetLogger()
	}
	return jh.logger
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

	jh.logger.PrintKeyValue("Before modification", "mapToProcess", mapToProcess)
	mapToProcess["degree"] = "phd"
	jh.logger.PrintKeyValue("After adding a new key-value", "mapToProcess", mapToProcess)

	return mapToProcess, nil
}
