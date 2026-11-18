package sources

import (
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/services/sources"
	"net/http"

	"github.com/google/uuid"
)

// GetSource
// @Tags         sources
// @Summary      Get a source.
// @Description  Get one feed.
// @Produce      json
// @Param        id             path      string  true  "Source id (uuid)"
// @Success      200            {object}  models.Source
// @Failure      404            {object}  models.Problem  "Source not found"
// @Failure      422            {object}  models.Problem  "A field is invalid"
// @Failure      500            {object}  models.Problem  "Something went wrong"
// @Router       /sources/{id} [get]
func GetSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sourceId, _ := ctx.Value("sourceId").(uuid.UUID) // getting key set in context.go

	source, err := sources.GetSourceById(&sourceId)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, source)
}
