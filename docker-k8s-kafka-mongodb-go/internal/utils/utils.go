package utils

import (
	"encoding/json"
	"fmt"
)

func StringToStruct(s string, o any) error {
	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		fmt.Println(err)
	}

	return err
}
