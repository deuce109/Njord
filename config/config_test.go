package config

import "testing"

func TestReadConfig(t *testing.T) {
	err := ReadConfig("../test_data/world.toml")
	if err != nil {
		t.Fatal(err)
	}

	if WorldName != "Test" {
		t.Errorf("Expected WorldName to be Test got, %s", WorldName)
	}

	if SaveName != "Test" {
		t.Errorf("Expected SaveName to be Test got, %s", SaveName)
	}
}

func TestReadConfigWithBadPath(t *testing.T) {
	err := ReadConfig("./")
	if err == nil {
		t.Fatal(err)
	}
}
