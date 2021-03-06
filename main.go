package main

import (
	"encoding/json"
	"fmt"
	"go-api-playground/cmd"
	"go-api-playground/logging"
	"go-api-playground/plugin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

var plugins plugin.Plugins
var conf map[string]interface{}
var logChannel chan string
var plugUpdateChannel chan string
var logger logging.Logger
var router *mux.Router

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	confFilePath := path.Join(dir, "conf.json")

	confJSON, err := os.Open(confFilePath)
	handleErr(err)
	defer confJSON.Close()

	confBytes, _ := ioutil.ReadAll(confJSON)

	json.Unmarshal([]byte(confBytes), &conf)

	logChannel = make(chan string, 1024)
	logger = logging.NewLogger(conf["log_file"].(string), logChannel)
	// go processLogEntries()
	go logger.ProcessLogEntries()

	var errs []error
	plugins, errs = plugin.GetPlugins(conf["route_directory"].(string))
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Log(fmt.Sprintf("error: unable to process one of the included plugins. %s", err.Error()))
			log.Fatal(err)
		}
	}

	plugUpdateChannel = make(chan string, 1024)

	go plugin.CreatePluginWatcher(conf["route_directory"].(string), plugUpdateChannel, logChannel)
	go processPluginUpdates()

	router = mux.NewRouter()

	for _, plg := range plugins.Plugins {
		logger.Log(fmt.Sprintf("info: processed route plugin [Type: %s; Path: %s; Method: %s;]", plg.Type, plg.Path, plg.Method))

		if plg.Type == "go" {
			router.HandleFunc(plg.Path, ExecGoRoute).Methods(plg.Method)
		} else if plg.Type == "python" {
			router.HandleFunc(plg.Path, ExecPyRoute).Methods(plg.Method)
		} else if plg.Type == "powershell" {
			router.HandleFunc(plg.Path, ExecPsRoute).Methods(plg.Method)
		}
	}

	logger.Log(fmt.Sprintf("info: starting server [Port: %s]", conf["port"].(string)))
	log.Fatal(http.ListenAndServe(conf["port"].(string), router))
	logger.Log("info: stopping server")

	logChannel <- "<EOC>"
}

// handleErr checks if err is nil and calls a fatal log if it is not.
// Arguments:
//   err is the error to check if nil
// Returns:
func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func processPluginUpdates() {
	for {
		upd := <-plugUpdateChannel

		if upd == "<EOC>" {
			break
		}
		if len(upd) > 0 {
			var errs []error
			plugins, errs = plugin.GetPlugins(conf["route_directory"].(string))
			if len(errs) > 0 {
				for _, err := range errs {
					logger.Log(fmt.Sprintf("error: unable to process one of the included plugins. %s", err.Error()))
					log.Fatal(err)
				}
			}

			for _, plg := range plugins.Plugins {
				logger.Log(fmt.Sprintf("info: processed route plugin [Type: %s; Path: %s; Method: %s;]", plg.Type, plg.Path, plg.Method))

				if plg.Type == "go" {
					router.HandleFunc(plg.Path, ExecGoRoute).Methods(plg.Method)
				} else if plg.Type == "python" {
					router.HandleFunc(plg.Path, ExecPyRoute).Methods(plg.Method)
				} else if plg.Type == "powershell" {
					router.HandleFunc(plg.Path, ExecPsRoute).Methods(plg.Method)
				}
			}
		}
	}
}

// ExecPyRoute executes a python route plugin.
// Arguments:
//   w is the ResponseWriter
//   r is the pointer to the Request
// Returns:
func ExecPyRoute(w http.ResponseWriter, r *http.Request) {
	var runVars []string

	params := mux.Vars(r)

	logger.Log(fmt.Sprintf("request: [RouteType: python; URL: %s; Method: %s; RemoteAddress: %s]", r.URL.Path, r.Method, r.RemoteAddr))

	var actPlug plugin.Plugin

	mstr := "^/" + strings.Split(r.URL.Path, "/")[1] + "($|/|\\b)"

	regex, err := regexp.Compile(mstr)
	handleErr(err)

	for _, plg := range plugins.Plugins {
		matched := regex.MatchString(plg.Path)
		if matched == true {
			actPlug = plg
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if len(actPlug.Path) > 0 {
		for _, param := range actPlug.Parameters {
			if len(params[param]) > 0 {
				runVars = append(runVars, params[param])
			}
		}

		pcmd := cmd.PyCommand{
			FilePath: fmt.Sprintf("%s/%s", conf["route_directory"], actPlug.MainFile),
			RunArgs:  runVars,
			Logger:   logging.NewLogger(conf["log_file"].(string), logChannel),
		}
		pcmd.Init()
		pcmd.Exec()

		fmt.Fprintf(w, pcmd.StdOut)
	} else {
		fmt.Fprintf(w, "{\"error\": \"plugin not found\"}")
	}
}

// ExecGoRoute executes a go route plugin.
// Arguments:
//   w is the ResponseWriter
//   r is the pointer to the Request
// Returns:
func ExecGoRoute(w http.ResponseWriter, r *http.Request) {
	var runVars []string

	params := mux.Vars(r)

	logger.Log(fmt.Sprintf("request: [RouteType: go; URL: %s; Method: %s; RemoteAddress: %s]", r.URL.Path, r.Method, r.RemoteAddr))

	var actPlug plugin.Plugin

	mstr := "^/" + strings.Split(r.URL.Path, "/")[1] + "($|/|\\b)"

	regex, err := regexp.Compile(mstr)
	handleErr(err)

	for _, plg := range plugins.Plugins {
		matched := regex.MatchString(plg.Path)
		if matched == true {
			actPlug = plg
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if len(actPlug.Path) > 0 {
		for _, param := range actPlug.Parameters {
			if len(params[param]) > 0 {
				runVars = append(runVars, params[param])
			}
		}

		gcmd := cmd.GoCommand{
			FilePath: fmt.Sprintf("%s/%s", conf["route_directory"], actPlug.MainFile),
			RunArgs:  runVars,
			Logger:   logging.NewLogger(conf["log_file"].(string), logChannel),
		}
		gcmd.Init()
		gcmd.Exec()

		fmt.Fprintf(w, gcmd.StdOut)
	} else {
		fmt.Fprintf(w, "{\"error\": \"plugin not found\"}")
	}
}

// ExecPsRoute executes a powershell route plugin.
// Arguments:
//   w is the ResponseWriter
//   r is the pointer to the Request
// Returns:
func ExecPsRoute(w http.ResponseWriter, r *http.Request) {
	var runVars []string

	params := mux.Vars(r)

	logger.Log(fmt.Sprintf("request: [RouteType: powershell; URL: %s; Method: %s; RemoteAddress: %s]", r.URL.Path, r.Method, r.RemoteAddr))

	var actPlug plugin.Plugin

	mstr := "^/" + strings.Split(r.URL.Path, "/")[1] + "($|/|\\b)"

	regex, err := regexp.Compile(mstr)
	handleErr(err)

	for _, plg := range plugins.Plugins {
		matched := regex.MatchString(plg.Path)
		if matched == true {
			actPlug = plg
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if len(actPlug.Path) > 0 {
		for _, param := range actPlug.Parameters {
			if len(params[param]) > 0 {
				runVars = append(runVars, params[param])
			}
		}

		pcmd := cmd.PsCommand{
			FilePath: fmt.Sprintf("%s/%s", conf["route_directory"], actPlug.MainFile),
			RunArgs:  runVars,
			Logger:   logging.NewLogger(conf["log_file"].(string), logChannel),
		}
		pcmd.Init()
		pcmd.Exec()

		fmt.Fprintf(w, pcmd.StdOut)
	} else {
		fmt.Fprintf(w, "{\"error\": \"plugin not found\"}")
	}
}
