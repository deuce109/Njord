package readers

import (
	"testing"
)

const GOOD_DATA_PATH = "../test_data/good"
const BAD_DATA_PATH = "../test_data/bad"
const BAD_PATH = "./bad/path"
const ERROR_NOT_NIL = "Error should not be nil"

func TestGetPlayerData(t *testing.T) {
	err := GetUserInfo("../test_data/good/", "../test_data/good/")

	if err != nil {
		t.Fatal(err)
	}

	if len(PlayerData) != 1 {
		t.Fatalf("Expected 1 players, got %d\n with error %s", len(PlayerData), err)
	}
}

func TestGetPlayerDataFailsOnBadAdminPath(t *testing.T) {
	err := GetUserInfo(GOOD_DATA_PATH, BAD_PATH)

	if err == nil {
		t.Fatal(ERROR_NOT_NIL)
	}

}

func TestGetPlayerDataFailsOnBadPlayerPath(t *testing.T) {
	err := GetUserInfo(BAD_PATH, GOOD_DATA_PATH)

	if err == nil {
		t.Fatal(ERROR_NOT_NIL)
	}

}

func TestGetPlayerDataFailsOnBadAdminXml(t *testing.T) {
	err := GetUserInfo(GOOD_DATA_PATH, BAD_DATA_PATH)

	if err == nil {
		t.Fatal(ERROR_NOT_NIL)
	}

}

func TestGetPlayerDataFailsOnBadPlayerXml(t *testing.T) {
	err := GetUserInfo(BAD_DATA_PATH, GOOD_DATA_PATH)

	if err == nil {
		t.Fatal(ERROR_NOT_NIL)
	}

}

func TestGetPlayerDataFailsPoorlySetPermissionLevels(t *testing.T) {
	err := GetUserInfo(GOOD_DATA_PATH, "../test_data/bad/atoi")

	if err == nil {
		t.Fatal(ERROR_NOT_NIL)
	}

}
