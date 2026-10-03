package sources

import (
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/services/sources"
	"net/http"
)

// GetSources
// @Tags         sources
// @Summary      Get all sources.
// @Description  Get every feed the ingestor knows.
// @Produce      json
// @Success      200            {array}   models.Source
// @Failure      500            {object}  models.Problem  "Something went wrong"
// @Router       /sources [get]
func GetSources(w http.ResponseWriter, _ *http.Request) {
	allSources, err := sources.GetAllSources()
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, allSources)
}
