package comment

// Repository представляет репозиторий для реализации.
type Repository interface {
	Creator
	Updater
	Finder
}
