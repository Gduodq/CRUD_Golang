package bookStruct

type Book struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Subtitle    string  `json:"subtitle"`
	AuthorName  string  `json:"authorName"`
	ReleaseDate string  `json:"releaseDate"`
	Price       float32 `json:"price"`
	CreatedAt   string  `json:"createdAt"`
}
