package logging

import (
	"fmt"
	"time"
)

// Logger is a struct used to push channels to a message.
type Logger struct {
	OutputFile string
	LogChannel chan string
}

// NewLogger creates a new Logger object and returns it.
// --------------------------------------------------------------
// Arguments:
//   outFile is the output file that logs are written to.
//   logChan is the channel to push strings (messages) to.
// --------------------------------------------------------------
// Returns:
//   Logger object.
// --------------------------------------------------------------
func NewLogger(outFile string, logChan chan string) Logger {
	return Logger{
		OutputFile: outFile,
		LogChannel: logChan,
	}
}

// Log extends the Logger object and pushes a formatted string (message) to the channel.
// --------------------------------------------------------------
// Arguments:
//   message is the string to format and push to the channel
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
func (l *Logger) Log(message string) {
	t := time.Now()
	l.LogChannel <- fmt.Sprintf("[%s]\t%s", t.Format("2006-01-02 15:04:05"), message)
}
