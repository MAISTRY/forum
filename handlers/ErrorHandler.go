package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorPage writes an HTTP error response with the specified error message, status code, and page name.
// The response is encoded as a JSON object with the following fields:
//   - PageName: the name of the page that encountered the error
//   - ErrorMessage: the error message to be displayed
//   - Statuscode: the HTTP status code of the error
func ErrorPage(w http.ResponseWriter, errMsg string, errCode int, pageName string) {
	w.WriteHeader(errCode)

	data := map[string]interface{}{
		"PageName":     pageName,
		"ErrorMessage": errMsg,
		"Statuscode":   pageName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
