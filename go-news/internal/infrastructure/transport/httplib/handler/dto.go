package handler

import uc "github.com/ee-crocush/go-news/go-news/internal/usecase/post"

// PostDTO представляет пост в массиве постов.
type PostDTO struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime string `json:"pub_time"`
}

func MapPostToPostDTO(post uc.PostDTO) PostDTO {
	return PostDTO{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Link:    post.Link,
		PubTime: post.PubTime,
	}
}

func MapNewsToNewsDTO(news []uc.PostDTO) []PostDTO {
	newsDTO := make([]PostDTO, 0, len(news))
	for _, post := range news {
		postDTO := MapPostToPostDTO(post)
		newsDTO = append(newsDTO, postDTO)
	}

	return newsDTO
}
