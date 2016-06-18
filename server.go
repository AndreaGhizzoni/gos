// web-server test
package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// application home directory
var appHome = "/home/andrea/.gos/"

// Variable to hold all template file just specify the patter where to find all
// the templates
var templates = template.Must(template.ParseGlob(appHome + "template/*.html"))

// default application log path
var defLogPath = appHome + "gos.log"

// function to initialize the logger
func initLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	flag := log.Ldate | log.Ltime | log.Lshortfile

	Trace = log.New(traceHandle, "TRACE: ", flag)
	Info = log.New(infoHandle, "INFO: ", flag)
	Warning = log.New(warningHandle, "WARNING: ", flag)
	Error = log.New(errorHandle, "ERROR: ", flag)
}

// function to check if home directory exists, if not create a new one
func initHomeDir() {
	if e, _ := exists(appHome); !e {
		err := os.Mkdir(appHome, 0700)
		if err != nil {
			panic(err)
		}
	}
}

// Auxiliary function to render a specific template.
// w writer where to write the template
// tmpl template name
// p page structure where to find the data
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Http handler for index
func rootHandler(w http.ResponseWriter, r *http.Request) {
	Trace.Print("Index handler called")
	renderTemplate(w, "index")
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	Trace.Print("Shutdown handler called")
	err := exec.Command("/sbin/poweroff").Run()
	if err != nil {
		Error.Fatal(err)
	}
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {
	initHomeDir()

	logFile, err := os.OpenFile(defLogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	if err != nil {
		log.Fatalln("Failed to open log file :", err)
	}
	defer func() {
		Trace.Print("Closing log file")
		logFile.Close()
	}()

	initLogger(logFile, logFile, logFile, logFile)

	Trace.Print("Program Start")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/shutdown", shutdownHandler)
	Trace.Print("Server Started")
	Error.Fatal(
		http.ListenAndServe(":8080", nil))
}
