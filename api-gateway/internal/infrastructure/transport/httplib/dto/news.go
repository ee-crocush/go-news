package dto

// FindByIDResponse Представляет структуру ответа для поиска по ID.
// Здесь вынесли отдельно в структуру, так как мы собираем данные из нескольких сервисов.
type FindByIDResponse struct {
	Data PostWithComments `json:"data"`
}

// PostWithComments описывает структуру новости с комментариями.
type PostWithComments struct {
	Post     Post      `json:"post"`
	Comments []Comment `json:"comments"`
}

// Post описывает структуру новости.
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime string `json:"pub_time"`
}
