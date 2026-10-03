package articles

import (
	"context"
	"fmt"
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Context
/* This method is used to get ressource ID from url
*
* In REST, urls are formed like this : articles/{specific_collection_ressource_id}/...
* Here it is articles/{id}, where id is the UUID v5 of the article (contracts/identity.md).
 */
func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleId, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			helpers.RespondError(w, fmt.Errorf("cannot parse id (%s) as UUID: %w", chi.URLParam(r, "id"), models.ErrUnprocessableEntity))
			return
		}

		ctx := context.WithValue(r.Context(), "articleId", articleId) // We fill context with a Key-valued variable
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
