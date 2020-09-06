package upload

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Upload struct {
	Path     string
	Template string
}

func (u *Upload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/// uploads will have multiple parts:
	/// 1. The actual file being uploaded, passed as 'file'
	/// 2. META- headers, containing extra data
	///    META- headers will be parsed into json, META-foo ends up in .foo
	if r.Method == http.MethodPost {
		file, fileheader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "missing file", http.StatusBadRequest)
			return
		}

		meta := make(map[string]interface{})

		for name, value := range r.Header {
			if strings.HasPrefix(name, "Meta-") {
				key := strings.ToLower(name)[5:]
				meta[key] = value[0]
			}
		}

		//write meta to meta/ using config.FileFormat as template
		tmpl, err := template.New("").Parse(u.Template)

		if err != nil {
			// error
		}

		buf := &bytes.Buffer{}
		tmpl.Execute(buf, meta)

		filePath := buf.String()
		filePath = strings.Replace(filePath, "..", "", -1)
		path.Clean(filePath)

		finalPath := filepath.Join(u.Path, filePath)

		_, err = os.Stat(finalPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(finalPath, 0755)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		metaContents, err := json.Marshal(meta)
		log.Println(string(metaContents))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(filepath.Join(finalPath, fileheader.Filename+".meta"), metaContents, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.OpenFile(filepath.Join(finalPath, fileheader.Filename), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(f, file)

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
