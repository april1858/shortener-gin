package entity

import "errors"

// Error codes returned when originalURL is missing.
var (
	ErrNotFound = errors.New("not found")
	ErrDeleted  = errors.New("deleted")
)

// Redirect is used to get a JSON response in the GetAllUID() function.
type Redirect struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ChData is used to transfer data over the channel for deletion.
type ChData struct {
	UID  string
	Data []string
}

// StoreElem is used to put the delete flag in memory (function funnelm() in memory.go)
type StoreElem struct {
	Short     string
	Original  string
	UID       string
	Condition bool
}
