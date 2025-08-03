package handler

// PostItem представляет пост в массиве постов.
type PostItem struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime string `json:"pub_time"`
}
