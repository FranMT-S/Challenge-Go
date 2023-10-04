package _logs

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	constants_log "github.com/FranMT-S/Challenge-Go/src/constants/logs"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
)

/*
LogSVG adds a new log record to the file. If the file does not exist it will create it

  - fileName the name of the log file

  - operation a name of the action where failed example database.

  - description a description personal

  - err any object errors detected
*/
func LogSVG(fileName, operation, description string, err error) {

	var file *os.File
	var header []string = nil

	Helpers.CreateDirectoryLogIfNotExist("logs")
	path := fmt.Sprintf("logs/%v.csv", fileName)

	if isNotExist(path) {
		header = []string{"Date", "Operation", "description", "Error"}
	}

	file, _err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if _err != nil {
		Println(constants_log.ERROR_CREATE_LOG)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err == nil {
		err = fmt.Errorf("")
	}

	line := []string{time.Now().Format(time.RFC1123), operation, description, err.Error()}

	if header != nil {
		writer.Write(header)
	}

	writer.Write(line)

	Println("\nrecord added to log file: " + path)

}

func isNotExist(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}

// Println Execute a log.Println but in red text
func Println(v ...any) {
	ColorRed()
	log.Println(v...)
	ColorWhite()
}

// Ansi Color

// Change color console print to red
func ColorRed() {
	fmt.Print("\033[31m")
}

// Change color console print to white
func ColorWhite() {
	fmt.Print("\033[0m")
}

// Change color console print to Green
func ColorGreen() {
	fmt.Print("\033[32m")
}

// Change color console print to Yellow
func ColorYellow() {
	fmt.Print("\033[33m")
}

// Change color console print to Blue
func ColorBlue() {
	fmt.Print("\033[34m")
}
