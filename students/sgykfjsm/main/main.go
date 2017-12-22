package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gophercises/urlshort/students/sgykfjsm"
)

func main() {
	yamlPtr := flag.String("yaml", "path_map.yaml", "the path to the path mapping yaml file.")
	isJsonPtr := flag.Bool("j", false, "If true, this program will use the specified path mapping json file.")
	jsonPtr := flag.String("json", "path_map.json", "the path to the path mapping json file.")
	flag.Parse()

	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var mappingFile string
	if *isJsonPtr {
		mappingFile = *jsonPtr
	} else {
		mappingFile = *yamlPtr
	}
	log.Printf("Use %s as mapping file", mappingFile)
	mappingData, err := ioutil.ReadFile(mappingFile)
	if err != nil {
		log.Fatal(err)
	}

	var handler http.HandlerFunc
	if *isJsonPtr {
		handler, err = urlshort.JSONHandler(mappingData, mapHandler)
	} else {
		handler, err = urlshort.YAMLHandler(mappingData, mapHandler)
	}
	if err != nil {
		panic(err)
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
