package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/lineageos-infra/go/pkg/auth"
	"github.com/lineageos-infra/go/pkg/mirror/files"
	"github.com/lineageos-infra/go/pkg/mirror/upload"
)

func main() {
	flag.String("path", "data", "Path to files on disk")
	flag.String("port", ":3000", "Port to listen on")
	flag.String("apikeys", "", "Commma seperated list of valid API keys")
	flag.Parse()

	log.Println("Serving data to/from ", flag.Lookup("path").Value.String())

	keys := make(map[string]bool)

	for _, key := range strings.Split(flag.Lookup("apikeys").Value.String(), ",") {
		keys[key] = true
	}

	fileService := files.Files{
		Path: flag.Lookup("path").Value.String(),
	}

	uploadService := upload.Upload{
		Path: flag.Lookup("path").Value.String(),
	}

	mux := http.NewServeMux()
	mux.Handle("/files", &fileService)
	mux.Handle("/upload", auth.ApiKey(&uploadService, keys))

	err := http.ListenAndServe(flag.Lookup("port").Value.String(), mux)

	if err != nil {
		log.Fatal(err)
	}
}
