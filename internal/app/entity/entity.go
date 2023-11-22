package entity

import "errors"

var ErrNotFound = errors.New("not found")
var ErrDeleted = errors.New("deleted")

type Redirect struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ChData struct {
	UID  string
	Data []string
}

type StoreElem struct {
	Short     string
	Original  string
	UID       string
	Condition bool
}
