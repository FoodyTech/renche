package apputil

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
)

// WriteJSONResponse encodes response into provided container.
func WriteJSONResponse(w http.ResponseWriter, code int, data interface{}) {
	raw, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(prettyJSON(raw))
}

func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", "  ")
	return out.Bytes()
}

// WriteHTMLResponse encodes response into provided container.
func WriteHTMLResponse(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(data)
}

// WriteZipResponse encodes response into provided container.
func WriteZipResponse(w http.ResponseWriter, code int, name string, data []byte) {
	// TODO: play with bufio

	buf := new(bytes.Buffer)
	wzip := zip.NewWriter(buf)

	f, _ := wzip.Create(name)
	_, _ = f.Write(data)
	wzip.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}
