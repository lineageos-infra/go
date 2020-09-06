package upload

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Upload struct {
	Path string
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
		log.Printf("%+v", meta)

		fn := fileheader.Filename

		f, err := os.OpenFile(filepath.Join(u.Path, fn), os.O_WRONLY|os.O_CREATE, 0644)
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
