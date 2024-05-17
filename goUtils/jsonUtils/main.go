package main

import "github/GoDevKit/goUtils/jsonUtils/json_utils"

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
	jsonHandler := json_utils.NewJsonHandler()
	jsonHandler.DisplayAllJsonHandlers(empStr)
}
