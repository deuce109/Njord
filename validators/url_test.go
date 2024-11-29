package validators

import (
	"testing"

	"gopkg.in/validator.v2"
)

type testUrl struct {
	Url string `validate:"path"`
}

type testInt struct {
	Int int `validate:"path"`
}

func TestMain(m *testing.M) {
	validator.SetValidationFunc("path", IsPath)
	m.Run()
}

func TestIsUrlForGoodUrl(t *testing.T) {
	test := &testUrl{
		Url: "https://www.google.com",
	}
	err := validator.Validate(test)

	if err != nil {
		t.Fatal(err)
	}
}

func TestIsUrlForFilepaths(t *testing.T) {
	test := &testUrl{
		Url: "./url.go",
	}

	err := validator.Validate(test)

	if err != nil {
		t.Fatal(err)
	}
}

func TestIsUrlOnEmptyString(t *testing.T) {
	test := &testUrl{
		Url: "",
	}

	err := validator.Validate(test)

	if err == nil {
		t.Fatal(err)
	}
}

func TestIsUrlFailsWhenGivenNonString(t *testing.T) {
	test := &testInt{
		Int: 123,
	}

	err := validator.Validate(test)

	if err == nil {
		t.Fatal(err)
	}
}
