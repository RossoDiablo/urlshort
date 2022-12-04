package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/RossoDiablo/urlshort"
)

const (
	confPath = "config.yaml"
)

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	mux := defaultMux()

	configFilePath := flag.String("path", confPath, "Use -path=filename to customize config filepath")
	flag.Parse()

	//reading file
	buf, err := os.ReadFile(*configFilePath)
	if err != nil {
		exit("Error opening file!")
	}

	ext := path.Ext(*configFilePath)

	var handler http.HandlerFunc
	switch {
	case ext == ".yaml":
		handler, err = urlshort.YAMLHandler(buf, mux)
		if err != nil {
			exit("Error handling YAML!")
		}
	case ext == ".json":
		handler, err = urlshort.JSONHandler(buf, mux)
		if err != nil {
			exit("Error handling JSON!")
		}
	default:
		exit("Incorrect format of config file!")
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
