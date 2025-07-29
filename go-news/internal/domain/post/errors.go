package post

import "errors"

var (
	// ErrInvalidPostID представляет ошибку невалидного идентификатора новости.
	ErrInvalidPostID = errors.New("invalid Post ID")
	// ErrEmptyPostTitle представляет ошибку незаполненного заголовка новости.
	ErrEmptyPostTitle = errors.New("empty Post title")
	// ErrEmptyPostContent представляет ошибку незаполненного содержимого новости.
	ErrEmptyPostContent = errors.New("empty Post content")
	// ErrEmptyPostLink представляет ошибку незаполненной ссылки источника новости.
	ErrEmptyPostLink = errors.New("empty Post link")
	// ErrEmptyPubTime представляет ошибку незаполненной даты публикации новости.
	ErrEmptyPubTime = errors.New("empty Publication time")
	// ErrPostNotFound представляет ошибку ненайденного поста.
	ErrPostNotFound = errors.New("post not found")
)
