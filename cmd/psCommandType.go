package cmd

import "go-api-playground/logging"

// PsCommand is the struct that is used with powershell plugins.
// It extends Command.
type PsCommand struct {
	Command
	RunArgs  []string `json:"runargs"`
	FilePath string   `json:"filepath"`

	Logger      logging.Logger
	LogChannerl chan string
}

// Init sets the Cmd to "pwsh", adds the "-File" argument, and appends any user provided args if any.
// --------------------------------------------------------------
// Arguments:
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
func (pc *PsCommand) Init() {
	pc.Cmd = "pwsh"
	pc.Arguments = []string{"-File", pc.FilePath}

	if len(pc.RunArgs) > 0 {
		for _, ra := range pc.RunArgs {
			pc.Arguments = append(pc.Arguments, ra)
		}
	}
}
