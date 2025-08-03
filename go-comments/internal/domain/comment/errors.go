package comment

import "errors"

var (
	// ErrInvalidCommentID представляет ошибку невалидного идентификатора комментария.
	ErrInvalidCommentID = errors.New("invalid comment ID")
	// ErrInvalidParentID представляет ошибку невалидного идентификатора родительского комментария.
	ErrInvalidParentID = errors.New("invalid parent ID")
	// ErrEmptyContent представляет ошибку незаполненного содержимого комментария.
	ErrEmptyContent = errors.New("empty comment content")
	// ErrWrongLengthUserName представляет ошибку по длине ника пользователя.
	ErrWrongLengthUserName = errors.New("username length must me between 6 and 50 symbols")
	// ErrEmptyPubTime представляет ошибку незаполненной даты публикации комментария.
	ErrEmptyPubTime = errors.New("empty Publication time")
	// ErrCommentNotFound представляет ошибку ненайденного комментария.
	ErrCommentNotFound = errors.New("comment not found")
)
