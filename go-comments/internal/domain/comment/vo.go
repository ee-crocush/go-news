package comment

import "time"

// ID - идентификатор комментария.
type ID struct {
	value int64
}

// NewID создает новый идентификатор комментария ID.
func NewID(id int64) (ID, error) {
	if id < 1 {
		return ID{}, ErrInvalidCommentID
	}
	return ID{value: id}, nil
}

// Value возвращает значение идентификатора комментария.
func (c ID) Value() int64 { return c.value }

// Equal сравнивает два идентификатора.
func (c ID) Equal(other ID) bool { return c.value == other.value }

// NewsID - идентификатор новости.
type NewsID struct {
	value int32
}

// Value возвращает значение идентификатора комментария.
func (c NewsID) Value() int32 { return c.value }

// NewNewsID создает новый идентификатор новости NewsID.
func NewNewsID(id int32) (NewsID, error) {
	if id < 1 {
		return NewsID{}, ErrInvalidNewsID
	}
	return NewsID{value: id}, nil
}

// ParentID - идентификатор родительского комментария.
type ParentID struct {
	value *int64
}

// NewParentID создает новый идентификатор родительского комментария ParentID.
func NewParentID(id int64) (ParentID, error) {
	if id < 1 {
		return ParentID{}, ErrInvalidParentID
	}
	return ParentID{value: &id}, nil
}

// NewEmptyParentID создаёт пустой (nil) ParentID.
func NewEmptyParentID() ParentID {
	return ParentID{value: nil}
}

// Value возвращает значение идентификатора родительского комментария.
func (p ParentID) Value() *int64 { return p.value }

// IsZero проверяет, установлен ли ParentID.
func (p ParentID) IsZero() bool {
	return p.value == nil
}

// Content - содержание комментария.
type Content struct {
	value string
}

// NewContent создает содержание комментария Content.
func NewContent(text string) (Content, error) {
	if len(text) > 0 {
		return Content{text}, nil
	}

	return Content{}, ErrEmptyContent
}

// Value возвращает значение содержания комментария.
func (c Content) Value() string { return c.value }

// CommentTime - время комментария.
type CommentTime struct {
	value time.Time
}

// NewTime создает время публикации комментария CommentTime.
func NewTime() CommentTime {
	return CommentTime{time.Now().UTC()}
}

// NewFromUnixSeconds создаёт CommentTime из секунд.
func NewFromUnixSeconds(s int64) (CommentTime, error) {
	if s <= 0 {
		return CommentTime{}, nil
	}

	return CommentTime{value: time.Unix(s, 0)}, nil
}

// Time возвращает значение времени публикации комментария.
func (t CommentTime) Time() time.Time {
	return t.value
}

// String возвращает строковое значение времени комментария в формате 2006-01-02 15:04:05
func (t CommentTime) String() string {
	return t.value.Format(time.DateTime)
}

// UserName - пользователь, оставивший комментарий.
type UserName struct {
	value string
}

// NewUserName создает пользователя UserName, оставившего комментарий.
func NewUserName(text string) (UserName, error) {
	if len(text) > 6 && len(text) < 50 {
		return UserName{text}, nil
	}

	return UserName{}, ErrWrongLengthUserName
}

// Value возвращает значение имени пользователя.
func (p UserName) Value() string { return p.value }

// Status представляет статус модерации комментария.
type Status struct {
	value string
}

const (
	// Pending модерация в процессе.
	Pending = "pending"
	// Approved Модерация пройдена.
	Approved = "approved"
	// Rejected Модерация не пройдена.
	Rejected = "rejected"
)

// NewStatus возвращает новый объект Status с заданным значением.
func NewStatus(status string) (Status, error) {
	switch status {
	case Pending, Approved, Rejected:
		return Status{value: status}, nil
	default:
		return Status{}, ErrInvalidStatus
	}
}

// Value возвращает значение статуса.
func (o Status) Value() string {
	return o.value
}
