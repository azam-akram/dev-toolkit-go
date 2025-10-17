package main

import (
	"dev-toolkit-go/utils-go/json-utils-go/logger"

	jsonhelper "github.com/azam-akram/json-helper-go"
)

var empStr = `{
    "id": "The ID",
    "name": "The User",
    "designation": "CEO",
    "address":
    [
        {
            "doorNumber": 1,
            "street": "The office street 1",
            "city": "The office city 1",
            "country": "The office country 1"
        },
        {
            "doorNumber": 2,
            "street": "The home street 2",
            "city": "The home city 2",
            "country": "The home country 2"
        }
    ]
}`

func main() {
	//jsonHandler := handler.NewJsonHandler()
	//jsonHandler.DisplayAllJsonHandlers(empStr)

	//fmt.Print("--------------------")

	jsonHelper := jsonhelper.NewJsonHelper()
	jsonHelper.StringToMap(empStr)

	logLevel := "INFO"
	log := logger.Init(logLevel)
	log.Info(empStr)

}
