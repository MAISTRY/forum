package handlers

import (
	"net/http"
	"strings"
)

// HomePage handles requests for the home page of the website.
// It performs URL validation, enforces the GET method, and serves the index.html file.
//
// Parameters:
//   - w: http.ResponseWriter - Used to construct an HTTP response.
//   - r: *http.Request - Contains information about the incoming HTTP request.
//
// This function does not return any value directly, but it may write to the
// http.ResponseWriter or perform HTTP redirects based on the request.
func HomePage(w http.ResponseWriter, r *http.Request) {

	segments := strings.Split(r.URL.Path, "/")

	if len(segments) > 3 && segments[1] != "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "./static/template/index.html")
}
