package entity

type Redirect struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ChData struct {
	UID  string
	Data []string
}
