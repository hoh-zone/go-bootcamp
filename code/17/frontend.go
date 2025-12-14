package chatserver

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/* web/assets/*
var embeddedFrontend embed.FS

// FrontendHandler serves the bundled static chat UI.
func FrontendHandler() http.Handler {
	content, err := fs.Sub(embeddedFrontend, "web")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(content))
}
