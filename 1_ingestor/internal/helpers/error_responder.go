package helpers

import (
	"encoding/json"
	"errors"
	"middleware/ingestor/internal/models"
	"net/http"

	"github.com/sirupsen/logrus"
)

// RespondError
// Writes the error response for err
func RespondError(w http.ResponseWriter, err error) {
	status := statusForError(err)

	// logging the real error, with its full context, for us
	logrus.WithError(err).Errorf("request failed with http code %d", status)

	var detail string
	// it's not a good practice to send the reason of an "internal server error" to the client
	if status == http.StatusInternalServerError {
		detail = "Something went wrong"
	} else {
		detail = err.Error()
	}

	writeProblem(w, status, detail)
}

// statusForError
// Chooses the HTTP status. errors.Is also finds errors wrapped with %w.
func statusForError(err error) int {
	switch {
	case errors.Is(err, models.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, models.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, models.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, models.ErrUnprocessableEntity):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func writeProblem(w http.ResponseWriter, status int, detail string) {
	problem := models.Problem{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	}
	body, _ := json.Marshal(problem)

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
