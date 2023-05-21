package request

import "encoding/json"

func JSON(data interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		return []byte("")
	}
	return []byte(result)
}
