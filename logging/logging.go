// Package logging contains the structs and methods for logging in a concurrent safe fashion.
package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger is a struct used to push channels to a message.
type Logger struct {
	OutputFile string
	LogChannel chan string
}

// NewLogger creates a new Logger object and returns it.
// Arguments:
//   outFile is the output file that logs are written to.
//   logChan is the channel to push strings (messages) to.
// Returns:
//   Logger object.
func NewLogger(outFile string, logChan chan string) Logger {
	return Logger{
		OutputFile: outFile,
		LogChannel: logChan,
	}
}

// Log extends the Logger object and pushes a formatted string (message) to the channel.
// Arguments:
//   message is the string to format and push to the channel
// Returns:
func (l *Logger) Log(message string) {
	t := time.Now()
	l.LogChannel <- fmt.Sprintf("[%s]\t%s", t.Format("2006-01-02 15:04:05"), message)
}

// ProcessLogEntries is run as a go routine and captures the messages
//   from the LogChannel and writes them to the log file.
// Arguments:
// Returns:
func (l *Logger) ProcessLogEntries() {
	lfile, err := os.Create(l.OutputFile)
	if err != nil {
		log.Fatal(err)
	}

	defer lfile.Close()

	for {
		msg := <-l.LogChannel

		if msg == "<EOC>" {
			break
		}

		fmt.Println(msg)
		lfile.WriteString(fmt.Sprintf("%s\n", msg))
	}
}
