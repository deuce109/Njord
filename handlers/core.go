package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func processBodyToItem[T interface{}](r *http.Request) (obj T, err error) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &obj)

	if err != nil {
		return
	}

	return
}
