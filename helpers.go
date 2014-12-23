package server

import (
	"html/template"
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