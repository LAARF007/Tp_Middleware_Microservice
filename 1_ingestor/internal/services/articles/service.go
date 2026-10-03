package articles

import (
	"database/sql"
	"errors"
	"fmt"
	"middleware/ingestor/internal/models"
	repository "middleware/ingestor/internal/repositories/articles"

	"github.com/google/uuid"
)

func GetLatestArticles() ([]models.IngestedArticle, error) {
	ingestedArticles, err := repository.GetAllArticles()
	if err != nil {
		return nil, fmt.Errorf("retrieving articles: %w", err)
	}

	return ingestedArticles, nil
}

func GetArticleById(id uuid.UUID) (*models.IngestedArticle, error) {
	article, err := repository.GetArticleById(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("article %s: %w", id, models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("retrieving article %s: %w", id, err)
	}
	return article, nil
}
