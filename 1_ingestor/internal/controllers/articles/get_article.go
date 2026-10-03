package articles

import (
	"github.com/google/uuid"
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/services/articles"
	"net/http"
)

// GetArticle
// @Tags         articles
// @Summary      Get an ingested article.
// @Description  Whether the ingestor has published this article, and from which source, when.
// @Produce      json
// @Param        id             path      string  true  "Article UUID (v5 of its canonical URL)"
// @Success      200            {object}  models.IngestedArticle
// @Failure      404            {object}  models.Problem  "Article not found"
// @Failure      422            {object}  models.Problem  "Cannot parse id"
// @Failure      500            {object}  models.Problem  "Something went wrong"
// @Router       /articles/{id} [get]
func GetArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articleId, _ := ctx.Value("articleId").(uuid.UUID) // getting key set in context.go

	ingestedArticle, err := articles.GetArticleById(articleId)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, ingestedArticle)
}
