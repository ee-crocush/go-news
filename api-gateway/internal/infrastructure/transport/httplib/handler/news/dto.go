package news

// PostResponse описывает структуру ответа на /news
type PostResponse struct {
	ID      int32  `json:"id" example:"1"`
	Title   string `json:"title" example:"Example title"`
	Content string `json:"content" example:"Example Long Content"`
	Link    string `json:"link" example:"https://example.com/news/1"`
	PubTime string `json:"pub_time" example:"2025-06-26 10:00:43"`
}
