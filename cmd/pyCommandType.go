package cmd

import "go-api-playground/logging"

// PyCommand is the struct that is used with python plugins.
// It extends Command.
type PyCommand struct {
	Command
	RunArgs  []string `json:"runargs"`
	FilePath string   `json:"filepath"`

	Logger      logging.Logger
	LogChannerl chan string
}

// Init sets the Cmd to "python" and appends any user provided args if any.
// --------------------------------------------------------------
// Arguments:
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
func (pc *PyCommand) Init() {
	pc.Cmd = "python"
	pc.Arguments = []string{pc.FilePath}

	if len(pc.RunArgs) > 0 {
		for _, ra := range pc.RunArgs {
			pc.Arguments = append(pc.Arguments, ra)
		}
	}
}
