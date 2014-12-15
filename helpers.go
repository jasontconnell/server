package server

import (
	"html/template"
)

func ToHtml(str string) template.HTML {
	return template.HTML(str)
}