// Package cmd contains the structs and methods for executing API route plugins.
package cmd

import (
	"errors"
	"os/exec"
)

// Command is the base struct for the plugin routes that is then extended to the account for different plugin languages.
type Command struct {
	Cmd       string   `json:"command"`
	Arguments []string `json:"argumets"`
	StdOut    string   `json:"stdout"`
	Err       string   `json:"err"`
}

// Exec executes the plugin.
// Arguments:
// Returns:
//   string holding the JSON from the plugin
//   error that occurred during plugin execution, if any.
func (cmd *Command) Exec() (string, error) {
	errString := ""

	if len(cmd.Cmd) <= 0 {
		if len(errString) <= 0 {
			errString += "error: the following properties must be set ["
		} else {
			errString += ", "
		}

		errString += "Cmd"
	}

	if len(cmd.Arguments) <= 0 {
		if len(errString) <= 0 {
			errString += "error: the following properties must be set ["
		} else {
			errString += ", "
		}

		errString += "Arguments"
	}

	if len(errString) > 0 {
		cmd.Err = errString
		return "", errors.New(errString)
	}

	out, err := exec.Command(cmd.Cmd, cmd.Arguments...).Output()
	if err != nil {
		cmd.Err = err.Error()
		return "", err
	}

	cmd.StdOut = string(out)

	return string(out), nil
}
