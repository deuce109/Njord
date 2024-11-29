package logging

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type LogOptions struct {
	Level          int
	OutputFolder   string
	EnableRotation bool
	WriteToConsole bool
}

const (
	TRACE int = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var opt LogOptions

func InitFromEnv(override ...string) {

	if os.Getenv("ENVRIONMENT") == "PRODUCTION" || (len(override) > 0 && override[0] == "PRODUCTION") {

		Init(LogOptions{
			EnableRotation: true,
			OutputFolder:   "/logs",
			Level:          INFO,
			WriteToConsole: true,
		})
	} else {

		Init(LogOptions{
			EnableRotation: false,
			Level:          TRACE,
			WriteToConsole: true,
		})
	}
}

func Init(options LogOptions) {
	opt = options
}

func checkLogLevel(lvl int) bool {
	return opt.Level <= lvl
}

func getLogLevelName(level int) string {
	var name string
	switch level {
	case 0:
		name = "TRACE"
	case 1:
		name = "DEBUG"
	case 2:
		name = "INFO"
	case 3:
		name = "WARN"
	case 4:
		name = "ERROR"
	case 5:
		name = "FATAL"
	}

	return name
}

func formatLog(msg string, lvl int) (output string) {

	levelName := getLogLevelName(lvl)

	output = fmt.Sprintf("%s | %s | %s", levelName, time.Now().String(), msg)

	return
}

func processOutputPath() (string, error) {
	var outputFileName string
	if opt.OutputFolder != "" {
		if output, isHomeRelative := strings.CutPrefix(opt.OutputFolder, "~"); isHomeRelative {
			if usr, err := user.Current(); err != nil {
				return "", err
			} else {
				outputFileName = path.Join(usr.HomeDir, output)
			}
		}

		if opt.EnableRotation {
			currentDate := time.Now()
			formattedDate := fmt.Sprintf("%d-%d-%d.log", currentDate.Year(), currentDate.Month(), currentDate.Day())
			outputFileName = filepath.Join(outputFileName, formattedDate)
		} else {
			outputFileName = filepath.Join(opt.OutputFolder, "server.log")
		}
	}

	return filepath.Abs(outputFileName)

}

func writeLog(msg string) {

	if outputFile, err := processOutputPath(); err != nil {
		fmt.Printf("%s", err)
	} else if outputFile != "" {
		if err = os.WriteFile(outputFile, []byte(msg), 0755); err != nil {
			fmt.Printf("%s", err)
		}
	}

	if opt.WriteToConsole {
		fmt.Println(msg)
	}
}

func checkErrorWrite(err error, lvl int) {
	if checkLogLevel(lvl) {
		writeLog(formatLog(err.Error(), lvl))
	}
}

func checkWrite(msg string, lvl int) {
	if checkLogLevel(lvl) {
		writeLog(formatLog(msg, lvl))
	}
}

func Trace(msg string) {
	checkWrite(msg, TRACE)
}

func Debug(msg string) {
	checkWrite(msg, DEBUG)
}

func Info(msg string) {
	checkWrite(msg, INFO)
}
func Warn(msg string) {
	checkWrite(msg, WARN)
}
func Error(msg string) {
	checkWrite(msg, ERROR)
}
func Fatal(msg string) {
	checkWrite(msg, FATAL)
}

func FatalError(err error) {
	checkErrorWrite(err, FATAL)
}
