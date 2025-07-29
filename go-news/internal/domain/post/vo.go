package post

import "time"

// PostID - идентификатор новости.
type PostID struct {
	value int32
}

// NewPostID создает новый идентификатор новости.
func NewPostID(id int32) (PostID, error) {
	if id < 1 {
		return PostID{}, ErrInvalidPostID
	}
	return PostID{value: id}, nil
}

// Value возвращает значение идентификатора новости.
func (t PostID) Value() int32 { return t.value }

// Equal сравнивает два идентификатора.
func (t PostID) Equal(other PostID) bool { return t.value == other.value }

// PostTitle - заголовок новости.
type PostTitle struct {
	value string
}

// NewPostTitle создает новый заголовок новости.
func NewPostTitle(text string) (PostTitle, error) {
	if len(text) > 0 {
		return PostTitle{text}, nil
	}

	return PostTitle{}, ErrEmptyPostTitle
}

// Value возвращает значение заголовка новости.
func (t PostTitle) Value() string { return t.value }

// PostContent - содержание новости.
type PostContent struct {
	value string
}

// NewPostContent создает содержание новости.
func NewPostContent(text string) (PostContent, error) {
	if len(text) > 0 {
		return PostContent{text}, nil
	}

	return PostContent{}, ErrEmptyPostContent
}

// Value возвращает значение содержания новости.
func (c PostContent) Value() string { return c.value }

// PubTime - время публикации новости.
type PubTime struct {
	value time.Time
}

// NewPubTime создает время публикации новости.
func NewPubTime() PubTime {
	return PubTime{time.Now().UTC()}
}

// NewFromUnixSeconds создаёт PubTime из секунд.
func NewFromUnixSeconds(s int64) (PubTime, error) {
	if s <= 0 {
		return PubTime{}, ErrEmptyPubTime
	}

	return PubTime{value: time.Unix(s, 0)}, nil
}

// Time возвращает значение времени публикации новости.
func (t PubTime) Time() time.Time {
	return t.value
}

// String возвращает строковое значение времени новости в формате 2006-01-02 15:04:05
func (t PubTime) String() string {
	return t.value.Format(time.DateTime)
}

// PostLink - ссылка на источник новости.
type PostLink struct {
	value string
}

// NewPostLink создает ссылку на источник новости.
func NewPostLink(text string) (PostLink, error) {
	if len(text) > 0 {
		return PostLink{text}, nil
	}

	return PostLink{}, ErrEmptyPostLink
}

// Value возвращает значение источника новости.
func (p PostLink) Value() string { return p.value }
