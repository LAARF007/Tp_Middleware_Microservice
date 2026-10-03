package sources

import (
	"database/sql"
	"middleware/ingestor/internal/helpers"
	"middleware/ingestor/internal/models"

	"github.com/google/uuid"
)

func GetAllSources() ([]models.Source, error) {
	rows, err := helpers.DB.Query("SELECT id, name, url FROM sources")
	if err != nil {
		return nil, err
	}
	// don't forget to close rows
	defer rows.Close()

	sources := []models.Source{}
	for rows.Next() {
		var source models.Source
		err = rows.Scan(&source.Id, &source.Name, &source.Url)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	return sources, rows.Err()
}

// GetSourceById
// Returns sql.ErrNoRows when there is no corresponding source.
func GetSourceById(id *uuid.UUID) (*models.Source, error) {
	rows, err := helpers.DB.Query("SELECT id, name, url FROM sources WHERE id = ?", id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	var source models.Source
	err = rows.Scan(&source.Id, &source.Name, &source.Url)
	if err != nil {
		return nil, err
	}

	return &source, nil
}

func InsertSource(source models.Source) error {
	_, err := helpers.DB.Exec(`INSERT INTO sources (id, name, url)
		VALUES (?, ?, ?)`,
		source.Id, source.Name, source.Url)
	return err
}
