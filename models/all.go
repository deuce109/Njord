package models

import "encoding/json"

func ToJson[T interface{}](t T) (data []byte, err error) {
	data, err = json.Marshal(&t)
	return
}

func ObjectFromJson[T interface{}](jsonData []byte) (t T, err error) {
	err = json.Unmarshal(jsonData, &t)
	return
}
