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
	path         string
	port         string
	apikeys      []string
	fileTemplate string
)

func init() {

	mirrorCmd.Flags().StringVar(&path, "path", "data", "Path to files on disk")
	mirrorCmd.Flags().StringVar(&port, "port", ":3000", "Port to listen on")
	mirrorCmd.Flags().StringSliceVar(&apikeys, "apikeys", nil, "List of api keys")
	mirrorCmd.Flags().StringVar(&fileTemplate, "template", "", "Template of file path to build from meta")

	viper.BindPFlag("path", mirrorCmd.Flags().Lookup("path"))
	viper.BindPFlag("port", mirrorCmd.Flags().Lookup("port"))
	viper.BindPFlag("apikeys", mirrorCmd.Flags().Lookup("apikeys"))

	rootCmd.AddCommand(mirrorCmd)
}

var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Starts service for listing/modifying/updating mirror files",
	Run:   mirror,
}

func mirror(cmd *cobra.Command, args []string) {

	keys := make(map[string]bool)

	for _, key := range viper.GetStringSlice("apikeys") {
		keys[key] = true
	}

	fileService := files.Files{
		Path: viper.GetString("path"),
	}

	uploadService := upload.Upload{
		Path:     viper.GetString("path"),
		Template: viper.GetString("template"),
	}

	mux := http.NewServeMux()
	mux.Handle("/files", &fileService)
	mux.Handle("/upload", auth.ApiKey(&uploadService, keys))

	err := http.ListenAndServe(viper.GetString("port"), mux)

	if err != nil {
		log.Fatal(err)
	}

}
