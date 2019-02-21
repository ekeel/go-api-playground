// Package cmd contains the structs and methods for executing API route plugins.
package cmd

import "go-api-playground/logging"

// GoCommand is the struct that is used with go plugins.
// It extends Command.
type GoCommand struct {
	Command
	RunArgs  []string `json:"runargs"`
	FilePath string   `json:"filepath"`

	Logger      logging.Logger
	LogChannerl chan string
}

// Init sets the Cmd to "go", add the "run" argument, and appends any user provided args if any.
// Arguments:
// Returns:
func (gc *GoCommand) Init() {
	gc.Cmd = "go"
	gc.Arguments = []string{"run", gc.FilePath}

	if len(gc.RunArgs) > 0 {
		gc.Arguments = append(gc.Arguments, "--")
		for _, ra := range gc.RunArgs {
			gc.Arguments = append(gc.Arguments, ra)
		}
	}
}
