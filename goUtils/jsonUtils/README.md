# [jsonUtils](https://github.com/azam-akram/utils-go/tree/main/json_utils):

Details about HTTP json-utils can be found [here](https://solutiontoolkit.com/2023/01/the-common-json-utility-functions-in-go-language/)

```
type JsonHandler_Interface interface {
	ConvertGenericInterfaceToMap()
	ConvertStringToMap(s string) (map[string]interface{}, error)
	ConvertMapToString(m map[string]interface{}) (string, error)
	ConvertStringToStruct(s string) (*Employee, error)
	ConvertStructToString(e *Employee) (string, error)
	ConvertStringToByte(s string) []byte
	ConvertByteToString([]byte) (string, error)
	ConvertByteToStruct(jsonBytes []byte) (*Employee, error)
	ConvertStructToByte(emp *Employee) (jsonBytes []byte, err error)
	ModifyInputJson(s string) (map[string]interface{}, error)
	DisplayAllJsonHandlers()
}
```
#### How to use:
Have a look at [DisplayAllJsonHandlers()](https://github.com/azam-akram/utils-go/blob/85de9b1f6804834765c9b0320d00ad944cac7b75/json_utils/json_handler.go#L54) to know how to use other functions exposed by `JsonHandler_Interface` interface. If you want to call all of these functions, simply call `DisplayAllJsonHandlers()`
```
jsonHandler := json_utils.JsonHandler{}
jsonHandler.DisplayAllJsonHandlers()
```