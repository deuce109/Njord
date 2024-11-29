package crypto

import (
	"crypto/rand"
	"testing"
)

var rng *Random = &Random{
	Reader: rand.Reader,
}

var secret string

func TestMain(m *testing.M) {
	secret, _ = rng.GetRandomSecret(32)
	m.Run()
}

const TEST_MSG = "Test message"

func TestHMACSHA512(t *testing.T) {

	sig := HMACSHA512(TEST_MSG, secret)

	if len(sig) != 44 {
		t.Fatalf("Expected `length` to be 44, but got %d", len(sig))
	}
}

func TestHMACEquals(t *testing.T) {

	sig1 := HMACSHA512(TEST_MSG, secret)
	sig2 := HMACSHA512(TEST_MSG, secret)

	if !HMACEquals(sig1, sig2) {
		t.Fatalf("Expected sig1 and sig2 to be equal \nSignature 1: %s\nSignature 2: %s", sig1, sig2)
	}
}
