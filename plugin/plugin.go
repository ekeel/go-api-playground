// Package plugin contains the code for handling API route plugins.
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
// Arguments:
//   pluginDir is the directory housing the plugin JSON/code files.
// Returns:
//   Plugins is the slice of found plugin(s)
//   []error is the slice of errors that occurred while loading the plugin(s)
func GetPlugins(pluginDir string) (Plugins, []error) {
	var errors []error
	var plugins Plugins

	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return plugins, append(errors, err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			fpath := path.Join(pluginDir, f.Name())

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

// // ProcessPluginUpdates monitors for changes being pushed to the plugin update channel and reloads plugins when they are changed.
// // Arguments:
// //   plugUpdateChannel is the channel to monitor for plugin updates.
// //   router is the mux Router that is hosting the site.
// // Returns:
// func ProcessPluginUpdates(plugUpdateChannel chan string, router *mux.Router) {
// 	for {
// 		upd := <-plugUpdateChannel

// 		if upd == "<EOC>" {
// 			break
// 		}
// 		if len(upd) > 0 {
// 			var errs []error
// 			plugins, errs := GetPlugins(conf["route_directory"].(string))
// 			if len(errs) > 0 {
// 				for _, err := range errs {
// 					logger.Log(fmt.Sprintf("error: unable to process one of the included plugins. %s", err.Error()))
// 					log.Fatal(err)
// 				}
// 			}

// 			for _, plg := range plugins.Plugins {
// 				logger.Log(fmt.Sprintf("info: processed route plugin [Type: %s; Path: %s; Method: %s;]", plg.Type, plg.Path, plg.Method))

// 				if plg.Type == "go" {
// 					router.HandleFunc(plg.Path, ExecGoRoute).Methods(plg.Method)
// 				} else if plg.Type == "python" {
// 					router.HandleFunc(plg.Path, ExecPyRoute).Methods(plg.Method)
// 				} else if plg.Type == "powershell" {
// 					router.HandleFunc(plg.Path, ExecPsRoute).Methods(plg.Method)
// 				}
// 			}
// 		}
// 	}
// }
