package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"plugsys/cmd"
	"plugsys/logging"
	"plugsys/plugin"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

var plugins plugin.Plugins
var conf map[string]interface{}
var logChannel chan string
var logger logging.Logger

func main() {
	confJSON, err := os.Open("/home/ekeeling/go/src/plugsys/conf.json")
	handleErr(err)
	defer confJSON.Close()

	confBytes, _ := ioutil.ReadAll(confJSON)

	json.Unmarshal([]byte(confBytes), &conf)

	logChannel = make(chan string, 1024)
	logger = logging.NewLogger(conf["log_file"].(string), logChannel)
	go processLogEntries()

	var errs []error
	plugins, errs = plugin.GetPlugins(conf["route_directory"].(string))
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Log(fmt.Sprintf("error: unable to process one of the included plugins. %s", err.Error()))
			log.Fatal(err)
		}
	}

	router := mux.NewRouter()

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
// --------------------------------------------------------------
// Arguments:
//   err is the error to check if nil
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// processLogEntries is run as a go routine and captures the messages
//   from the LogChannel and writes them to the log file.
// --------------------------------------------------------------
// Arguments:
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
func processLogEntries() {
	lfile, err := os.Create(logger.OutputFile)
	if err != nil {
		log.Fatal(err)
	}

	defer lfile.Close()

	for {
		msg := <-logChannel

		if msg == "<EOC>" {
			break
		}

		fmt.Println(msg)
		lfile.WriteString(fmt.Sprintf("%s\n", msg))
	}
}

// ExecPyRoute executes a python route plugin.
// --------------------------------------------------------------
// Arguments:
//   w is the ResponseWriter
//   r is the pointer to the Request
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
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
// --------------------------------------------------------------
// Arguments:
//   w is the ResponseWriter
//   r is the pointer to the Request
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
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
// --------------------------------------------------------------
// Arguments:
//   w is the ResponseWriter
//   r is the pointer to the Request
// --------------------------------------------------------------
// Returns:
// --------------------------------------------------------------
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
