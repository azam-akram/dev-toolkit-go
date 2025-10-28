package handler

type Handler interface {
	StringToStruct(string, any) error
	StructToString(any) (string, error)
	StringToMap(string) (map[string]any, error)
	MapToString(map[string]any) (string, error)
	BytesToString([]byte) string
	StringToBytes(string) []byte
	StructToBytes(any) ([]byte, error)
	BytesToStruct([]byte, any) error
	ModifyInputJson(string) (map[string]any, error)
	ProcessGenericMap(map[string]any)
	DisplayAllJsonHandlers(string)
}
