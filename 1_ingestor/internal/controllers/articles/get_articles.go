package articles

import (
	"fmt"
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/models"
	"middleware/ingestor/internal/services/articles"
	"net/http"

	"github.com/google/uuid"
)

// GetArticles
// @Tags         articles
// @Summary      Get the latest ingested articles.
// @Description  What the ingestor has published on ARTICLES.ingested, most recent first: id,
// @Description  urls, source and date only. The content travels in the NATS message.
// @Produce      json
// @Param        source_id      query     string  false  "Only the articles of this source (uuid)"
// @Param        limit          query     int     false  "How many articles; all of them when left out"
// @Success      200            {array}   models.IngestedArticle
// @Failure      422            {object}  models.Problem  "limit is not a number of 1 or more"
// @Failure      500            {object}  models.Problem  "Something went wrong"
// @Router       /articles [get]
func GetArticles(w http.ResponseWriter, r *http.Request) {
	// TODO get parameters
	sourceId := r.URL.Query().Get("source_id")
	_, err := uuid.Parse(sourceId)
	if err != nil {
		helpers.RespondError(w, fmt.Errorf("cannot parse source_id (%s) as UUID: %w", sourceId, models.ErrUnprocessableEntity))
		return
	}

	ingestedArticles, err := articles.GetLatestArticles()
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, ingestedArticles)
}
