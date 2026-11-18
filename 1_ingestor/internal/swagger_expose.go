package internal

// @title        middleware/ingestor
// @version      1.0.0
// @description  Ingestor of the technology-watch pipeline: polls RSS/Atom feeds and publishes
// @description  every new article on the NATS subject ARTICLES.ingested.
// @description  Errors follow RFC 9457 (application/problem+json), see contracts/http-errors.md.

// @schemes   http
// @BasePath  /
