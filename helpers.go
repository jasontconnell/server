package server

import (
	"html/template"
	"strings"
)

func ToHtml(str string) template.HTML {
	return template.HTML(str)
}

func ToKey(str string) template.HTML {
	return template.HTML(strings.ToLower(strings.Replace(str, " ", "-", -1)))
}