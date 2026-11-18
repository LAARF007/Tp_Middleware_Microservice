package sources

import (
	"database/sql"
	"errors"
	"fmt"
	"middleware/ingestor/internal/models"
	repository "middleware/ingestor/internal/repositories/sources"

	"github.com/google/uuid"
)

func GetAllSources() ([]models.Source, error) {
	sources, err := repository.GetAllSources()
	if err != nil {
		return nil, fmt.Errorf("retrieving sources: %w", err)
	}
	return sources, nil
}

func GetSourceById(id *uuid.UUID) (*models.Source, error) {
	source, err := repository.GetSourceById(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("source %s: %w", id, models.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("retrieving source %s: %w", id, err)
	}
	return source, nil
}

func CreateSource(source models.Source) (*models.Source, error) {
	err := validateSource(source)
	if err != nil {
		return nil, err
	}

	// It's possible to check whether the source already exists
	/*_, err = GetSourceByUrl(source.Url)
	if err == nil {
		return nil, fmt.Errorf("source %s already exists: %w", source.Id, models.ErrConflict)
	}
	if !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}*/

	newId := uuid.New()
	source.Id = &newId
	err = repository.InsertSource(source)
	if err != nil {
		return nil, fmt.Errorf("creating source %s: %w", source.Id, err)
	}
	return GetSourceById(source.Id)
}
