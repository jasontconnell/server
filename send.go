package server

import (
	"net/http"
	"io"
	"compress/gzip"
	"strings"
)


func sendContent (content []byte, w http.ResponseWriter, req *http.Request){
	writer := w.(io.Writer)
	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		writer = gzip.NewWriter(w)
		defer writer.(*gzip.Writer).Close()
	}
	writer.Write(content)
}