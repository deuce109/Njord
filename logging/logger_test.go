package logging

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()

	dirPath := "./"
	pattern := "*.log" // Example pattern: all text files

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file matches the pattern
		if matched, _ := filepath.Match(pattern, info.Name()); matched {
			// Delete the file
			err := os.Remove(path)
			if err != nil {
				return err
			}
			fmt.Println("Deleted:", path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

}

func TestInitFromEnvironment(t *testing.T) {
	InitFromEnv()
}

func TestInitFromEnvironmentWithProdSet(t *testing.T) {
	InitFromEnv("PRODUCTION")
}

func TestDebug(t *testing.T) {
	Debug("test")
}

func TestTrace(t *testing.T) {
	Trace("test")
}
func TestInfo(t *testing.T) {
	Info("test")
}
func TestWarnWithString(t *testing.T) {
	Warn("test")
}
func TestWarnWithError(t *testing.T) {
	Error("test")
}
func TestErrorWithError(t *testing.T) {
	Error("test")
}
func TestFatal(t *testing.T) {
	Fatal("test")
}

func TestFatalError(t *testing.T) {
	FatalError(errors.New("test"))
}

func TestInitWithoutRotation(t *testing.T) {
	Init(LogOptions{
		Level:          TRACE,
		OutputFolder:   "~/test_data/logs",
		WriteToConsole: true,
		EnableRotation: false,
	})

	Trace("Test")
}

func TestInitThatDisablesLogging(t *testing.T) {
	Init(LogOptions{
		Level:          TRACE,
		OutputFolder:   "",
		WriteToConsole: false,
		EnableRotation: false,
	})

	Trace("Test")
}
