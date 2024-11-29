package crypto

import (
	"crypto/rand"
	"errors"
	"math"
	"testing"
	"testing/iotest"
)

func testGetRandomSecretWithLength(t *testing.T, length int) {
	expected := (length + int(math.Ceil(float64(length)/3.0)))

	_rand := &Random{
		Reader: rand.Reader,
	}

	secret, err := _rand.GetRandomSecret(length)

	if err != nil {
		t.Fatal(err)
	} else if len(secret) != expected {
		t.Fatalf("len(secret) got %d want %d", len(secret), expected)
	}
}

func TestGetRandomSecret(t *testing.T) {
	testGetRandomSecretWithLength(t, 8)
	testGetRandomSecretWithLength(t, 16)
	testGetRandomSecretWithLength(t, 32)
	testGetRandomSecretWithLength(t, 64)
	testGetRandomSecretWithLength(t, 128)
}

func TestGetRandomSecretFailure(t *testing.T) {
	_rand := &Random{
		Reader: iotest.ErrReader(errors.New("Test for failure")),
	}

	_, err := _rand.GetRandomSecret(10)

	if err == nil {
		t.Fatal("Expected error to not be nil")
	}
}
