package server

import (
	"net/url"
)

type Page struct {
	Url url.URL
	Path string
}