package server

import (
	"path/filepath"
)

func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	var mime string

	switch ext {
		case ".js": mime = "application/javascript"
		case ".json": mime = "application/json"
		case ".css": mime = "text/css"
	}

	return mime
}