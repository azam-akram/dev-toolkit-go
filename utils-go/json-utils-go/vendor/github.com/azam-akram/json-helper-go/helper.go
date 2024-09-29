package jsonhelper

type Helper interface {
	StringToStruct(s string, i interface{}) error
	StructToString(i interface{}) (string, error)
	StringToMap(s string) (map[string]interface{}, error)
	MapToString(m map[string]interface{}) (string, error)
	BytesToString(jsonBytes []byte) string
	StringToBytes(s string) []byte
	StructToBytes(i interface{}) (jsonBytes []byte, err error)
	BytesToStruct(b []byte, d interface{}) error
	ModifyInputJson(s string) (map[string]interface{}, error)
}
