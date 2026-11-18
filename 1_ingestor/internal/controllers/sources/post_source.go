package sources

import (
	"encoding/json"
	"fmt"
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/models"
	"middleware/ingestor/internal/services/sources"
	"net/http"
)

// PostSource
// @Tags         sources
// @Summary      Create a source.
// @Description  Add a feed to poll.
// @Accept       json
// @Produce      json
// @Param        source         body      models.Source   true  "The source"
// @Success      201            {object}  models.Source
// @Header       201            {string}  Location        "/sources/{id}"
// @Failure      400            {object}  models.Problem  "Body is not valid JSON"
// @Failure      409            {object}  models.Problem  "A source with this id already exists"
// @Failure      422            {object}  models.Problem  "A field is invalid"
// @Failure      500            {object}  models.Problem  "Something went wrong"
// @Router       /sources [post]
func PostSource(w http.ResponseWriter, r *http.Request) {
	var source models.Source
	err := json.NewDecoder(r.Body).Decode(&source)
	if err != nil {
		helpers.RespondError(w, fmt.Errorf("body is not a valid source JSON object (%s): %w", err.Error(), models.ErrBadRequest))
		return
	}

	createdSource, err := sources.CreateSource(source)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, createdSource)
}
