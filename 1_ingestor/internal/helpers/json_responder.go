package helpers

import (
	"encoding/json"
	"net/http"
)

// RespondJSON
// Writes a successful JSON response: the status, the Content-Type header, then the body.
func RespondJSON(w http.ResponseWriter, status int, value any) {
	body, err := json.Marshal(value)
	if err != nil {
		RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
