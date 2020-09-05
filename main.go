package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/lineageos-infra/mirrorservice/pkg/files"
)

func main() {
	flag.String("path", "data", "Path to files on disk")
	flag.String("port", ":3000", "Port to listen on")
	flag.Parse()

	log.Println("Serving data to/from ", flag.Lookup("path").Value.String())

	fileService := files.Files{
		Path: flag.Lookup("path").Value.String(),
	}

	mux := http.NewServeMux()
	mux.Handle("/files", &fileService)

	err := http.ListenAndServe(flag.Lookup("port").Value.String(), mux)

	if err != nil {
		log.Fatal(err)
	}
}
