package files

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Files struct {
	Path string
}

func (u *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We store files on disk in the following format:
	// path/to/mirror/file.zip
	// path/to/mirror/file.zip.meta

	var metadata []json.RawMessage
	err := filepath.Walk(u.Path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".meta" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		metadata = append(metadata, json.RawMessage(data))
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return 500
	}
	w.Header().Set("Content-Type", "application/json")
	d, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		// 500
	}
	w.Write(d)
}
