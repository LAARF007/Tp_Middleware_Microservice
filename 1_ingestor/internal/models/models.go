package models

import (
	"time"

	"github.com/google/uuid"
)

type Source struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
	Url  string     `json:"url"`
}

type IngestedArticle struct {
	Id           uuid.UUID `json:"id"`
	Url          string    `json:"url"`
	CanonicalUrl string    `json:"canonical_url"`
	SourceId     string    `json:"source_id"`
	FirstSeenAt  time.Time `json:"first_seen_at" format:"date-time" example:"2026-09-13T15:07:02Z"`
}
