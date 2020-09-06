package cmd

import (
	"log"
	"net/http"

	"github.com/lineageos-infra/go/pkg/auth"
	"github.com/lineageos-infra/go/pkg/mirror/files"
	"github.com/lineageos-infra/go/pkg/mirror/upload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	path    string
	port    string
	apikeys []string
)

func init() {

	mirrorserviceCmd.Flags().StringVar(&path, "path", "data", "Path to files on disk")
	mirrorserviceCmd.Flags().StringVar(&port, "port", ":3000", "Port to listen on")
	mirrorserviceCmd.Flags().StringArrayVar(&apikeys, "apikeys", nil, "List of api keys")

	viper.BindPFlag("path", mirrorserviceCmd.Flags().Lookup("path"))
	viper.BindPFlag("port", mirrorserviceCmd.Flags().Lookup("port"))
	viper.BindPFlag("apikeys", mirrorserviceCmd.Flags().Lookup("apikeys"))

	rootCmd.AddCommand(mirrorserviceCmd)
}

var mirrorserviceCmd = &cobra.Command{
	Use:   "mirrorservice",
	Short: "Starts service for listing/modifying/updating mirror files",
	Run:   mirrorservice,
}

func mirrorservice(cmd *cobra.Command, args []string) {

	log.Println("Serving data to/from", path)

	keys := make(map[string]bool)

	for _, key := range apikeys {
		keys[key] = true
	}

	fileService := files.Files{
		Path: path,
	}

	uploadService := upload.Upload{
		Path: path,
	}

	mux := http.NewServeMux()
	mux.Handle("/files", &fileService)
	mux.Handle("/upload", auth.ApiKey(&uploadService, keys))

	err := http.ListenAndServe(port, mux)

	if err != nil {
		log.Fatal(err)
	}

}
