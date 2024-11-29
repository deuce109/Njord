package repositories

import (
	"testing"
)

func TestGetPlayerInfoReturnsValidUser(t *testing.T) {
	username := "testUserName"

	p := GetPlayerInfo(username)

	if p == nil {
		t.Error("Expected user to be returned")
	}
}

func TestGetPlayerInfoReturnsNilWithInvalidUserName(t *testing.T) {
	username := ""

	p := GetPlayerInfo(username)

	if p != nil {
		t.Error("Did not expect user to be returned")
	}
}
