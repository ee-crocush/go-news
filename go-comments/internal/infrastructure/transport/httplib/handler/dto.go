package handler

// CommentItem представляет комментарий в массиве комментариев.
type CommentItem struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime string `json:"pub_time"`
}
