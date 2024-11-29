package readers

import "testing"

func TestReadSecretFileFailsWithNoFileToRead(t *testing.T) {
	_, err := ReadSecretFile("./test_data/secrets/bad")
	if err == nil {
		t.Fatal("Error should not be nil")
	}
}

func TestReadSecretFile(t *testing.T) {
	_, err := ReadSecretFile("../test_data/good/secret")
	if err != nil {
		t.Fatal(err)
	}
}
