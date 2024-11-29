package models

import (
	"encoding/json"
	"testing"
)

type testObj struct {
	Name string `json:"name"`
}

var test = &testObj{
	Name: "test",
}

func TestToJson(t *testing.T) {
	jsonData, err := ToJson[testObj](*test)

	if err != nil {
		t.Fatal("Object did not parse into json correctly\n", err, "\n")
	} else {
		var js json.RawMessage

		err = json.Unmarshal(jsonData, &js)
		if err != nil {
			t.Fatal("Json string did not parse into object correctly\n", err, "\n")
		}
	}
}

func TestFromJson(t *testing.T) {
	jsonData := []byte(`{
		"Name": "test"
	}`)

	_, err := ObjectFromJson[testObj](jsonData)

	if err != nil {
		t.Fatal("Json string did not parse into object correctly\n", err, "\n")
	}
}
