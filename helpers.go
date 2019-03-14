package server

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

func ToHtml(str string) template.HTML {
	return template.HTML(str)
}

func ToKey(str string) template.HTML {
	return template.HTML(strings.ToLower(strings.Replace(str, " ", "-", -1)))
}

func FormatDate(date time.Time, format string) string {
	return date.Format(format)
}

func SendJson(w http.ResponseWriter, payload interface{}) {
	encoder := json.NewEncoder(w)
	encoder.Encode(payload)
}

func DecodeJson(reader io.Reader, payload interface{}) {
	decoder := json.NewDecoder(reader)
	decoder.Decode(payload)
}

func SendError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}
