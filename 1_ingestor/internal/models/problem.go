package models

// Problem is the body of every error response.
// It follows RFC 9457 "Problem Details for HTTP APIs", and is sent with the
// header Content-Type: application/problem+json
type Problem struct {
	Type   string `json:"type"`             // always "about:blank" in this TP
	Title  string `json:"title"`            // standard phrase of the status, e.g. "Not Found"
	Status int    `json:"status"`           // same number as the HTTP status code
	Detail string `json:"detail,omitempty"` // what went wrong
}
