package bookStruct

type Book struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	Title       string  `json:"title" bson:"title,omitempty"`
	AuthorName  string  `json:"authorName" bson:"authorName,omitempty"`
	ReleaseDate string  `json:"releaseDate" bson:"releaseDate,omitempty"`
	Price       float32 `json:"price" bson:"price,omitempty"`
}
