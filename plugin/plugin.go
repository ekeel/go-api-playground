package plugin

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Plugin holds the unmarshaled plugin data from the plugins JSON file.
type Plugin struct {
	Path       string   `json:"path"`
	MainFile   string   `json:"main_file"`
	Method     string   `json:"method"`
	Type       string   `json:"type"`
	Parameters []string `json:"parameters"`
}

// Plugins holds a list of Plugin(s)
type Plugins struct {
	Plugins []Plugin
}

// GetPlugins find all of the plugins in the plugin directory.
// --------------------------------------------------------------
// Arguments:
//   pluginDir is the directory housing the plugin JSON/code files.
// --------------------------------------------------------------
// Returns:
//   Plugins is the slice of found plugin(s)
//   []error is the slice of errors that occurred while loading the plugin(s)
// --------------------------------------------------------------
func GetPlugins(pluginDir string) (Plugins, []error) {
	var errors []error
	var plugins Plugins

	files, err := ioutil.ReadDir("/home/ekeeling/go/src/plugsys/routes")
	if err != nil {
		return plugins, append(errors, err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			fpath := path.Join("/home/ekeeling/go/src/plugsys/routes", f.Name())

			plugcJSON, err := os.Open(fpath)
			if err != nil {
				errors = append(errors, err)
			} else {
				plugcBytes, err := ioutil.ReadAll(plugcJSON)
				if err != nil {
					errors = append(errors, err)
				}

				var plug Plugin
				json.Unmarshal(plugcBytes, &plug)

				plugins.Plugins = append(plugins.Plugins, plug)
			}
		}
	}

	return plugins, errors
}
