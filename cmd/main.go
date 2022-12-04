package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/RossoDiablo/urlshort"
	"github.com/boltdb/bolt"
)

const confPath = "config.yaml"

var configFilePath = flag.String("path", confPath, "Flag for customizing config filepath")

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	mux := defaultMux()

	flag.Parse()

	ext := path.Ext(*configFilePath)

	var handler http.HandlerFunc
	switch {
	case ext == ".yaml":
		buf, err := os.ReadFile(*configFilePath)
		if err != nil {
			exit("Error opening yaml file!")
		}
		handler, err = urlshort.YAMLHandler(buf, mux)
		if err != nil {
			exit("Error handling YAML!")
		}
	case ext == ".json":
		buf, err := os.ReadFile(*configFilePath)
		if err != nil {
			exit("Error opening json file!")
		}
		handler, err = urlshort.JSONHandler(buf, mux)
		if err != nil {
			exit("Error handling JSON!")
		}
	case ext == ".db":
		db, err := bolt.Open(*configFilePath, 0600, nil)
		if err != nil {
			exit("Error opening DB file!")
		}
		defer db.Close()

		handler = urlshort.DBHandler(db, mux)
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
