package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/RossoDiablo/urlshort"
)

const (
	path = "config.yaml"
)

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	configFilePath := flag.String("path", path, "Use -path=filename to customize config filepath")
	flag.Parse()

	//reading file
	buf, err := os.ReadFile(*configFilePath)
	if err != nil {
		exit("Error opening file!")
	}

	yamlHandler, err := urlshort.YAMLHandler(buf, mapHandler)
	if err != nil {
		exit("Error parsing yaml!")
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
