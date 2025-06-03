package post

type Post struct {
	Id       int
	Category string `json:"category"`
	Type     string `json:"type"`
	Url      string `json:"url"`
	Text     string `json:"text"`
	Title    string `json:"title"`
}
