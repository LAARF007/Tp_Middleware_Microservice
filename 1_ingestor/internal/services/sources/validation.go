package sources

import (
	"fmt"
	"middleware/ingestor/internal/models"
	"net/url"
	"strings"
)

// Polling a feed more often than every 5 minutes might get you blocked.
const minimumPollIntervalS = 300

// validateSource
// Checks every field a client sends
func validateSource(source models.Source) error {
	if strings.TrimSpace(source.Name) == "" {
		return fmt.Errorf("name must not be empty: %w", models.ErrUnprocessableEntity)
	}

	_, err := url.ParseRequestURI(source.Url)
	if err != nil {
		return fmt.Errorf("url %q invalid: %w", source.Url, models.ErrUnprocessableEntity)
	}

	return nil
}
