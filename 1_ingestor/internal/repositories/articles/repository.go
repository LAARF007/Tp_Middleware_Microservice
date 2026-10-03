package articles

import (
	"database/sql"
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/models"

	"github.com/google/uuid"
)

func GetAllArticles() ([]models.IngestedArticle, error) {
	rows, err := helpers.DB.Query("SELECT id, url, canonical_url, source_id, first_seen_at FROM articles")
	if err != nil {
		return nil, err
	}
	return scanAllIngestedArticles(rows)
}

// GetArticleById
// Returns sql.ErrNoRows when there is no article.
func GetArticleById(id uuid.UUID) (*models.IngestedArticle, error) {
	rows, err := helpers.DB.Query("SELECT id, url, canonical_url, source_id, first_seen_at FROM articles WHERE id = ?", id.String())
	if err != nil {
		return nil, err
	}

	// don't forget to close rows
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	ingestedArticle, err := scanIngestedArticle(rows)
	if err != nil {
		return nil, err
	}
	return &ingestedArticle, nil
}

// scanAllIngestedArticles
// Reads every row, then closes them
func scanAllIngestedArticles(rows *sql.Rows) ([]models.IngestedArticle, error) {
	defer rows.Close()

	ingestedArticles := []models.IngestedArticle{}
	for rows.Next() {
		ingestedArticle, err := scanIngestedArticle(rows)
		if err != nil {
			return nil, err
		}
		ingestedArticles = append(ingestedArticles, ingestedArticle)
	}
	return ingestedArticles, rows.Err()
}

// scanIngestedArticle
// Reads the current row into an IngestedArticle.
func scanIngestedArticle(rows *sql.Rows) (models.IngestedArticle, error) {
	var ingestedArticle models.IngestedArticle
	err := rows.Scan(&ingestedArticle.Id, &ingestedArticle.Url, &ingestedArticle.CanonicalUrl,
		&ingestedArticle.SourceId, &ingestedArticle.FirstSeenAt)
	return ingestedArticle, err
}
