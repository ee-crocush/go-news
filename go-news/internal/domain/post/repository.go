package post

// Repository представляет репозиторий для реализации.
type Repository interface {
	PostStore
	PostFinder
}
