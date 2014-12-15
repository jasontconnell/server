package server

import (
	"net/http"
	"io"
	"compress/gzip"
	"strings"
	"fmt"
)


func sendContent (content []byte, w http.ResponseWriter, req *http.Request){
	writer := w.(io.Writer)
	fmt.Println("gzip testing")
	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		fmt.Println("gzipping")
		w.Header().Add("Content-Encoding", "gzip")
		writer = gzip.NewWriter(w)
		defer writer.(*gzip.Writer).Close()
	}
	writer.Write(content)
}