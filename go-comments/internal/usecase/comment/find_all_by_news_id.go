package comment

import (
	"context"
	"fmt"
	dom "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
	"sort"
)

var _ FindAllByNewsIDContract = (*FindAllByNewsIDUseCase)(nil)

// FindAllByNewsIDUseCase представляет структуру, реализующую бизнес-логику для поиска всех комментариев.
type FindAllByNewsIDUseCase struct {
	repo dom.Repository
}

// NewFindAllByNewsIDUseCase создает новый экземпляр usecase для поиска всех комментариев.
func NewFindAllByNewsIDUseCase(repo dom.Repository) *FindAllByNewsIDUseCase {
	return &FindAllByNewsIDUseCase{repo: repo}
}

// Execute выполняет бизнес-логику поиска всех комментариев.
func (uc *FindAllByNewsIDUseCase) Execute(ctx context.Context, in AllByNewsIDDTO) ([]CommentDTO, error) {
	newsID, err := dom.NewNewsID(in.NewsID)
	if err != nil {
		return []CommentDTO{}, fmt.Errorf("FindAllByNewsIDUseCase.Execute: %w", err)
	}

	comments, err := uc.repo.FindAllByNewsID(ctx, newsID)
	if err != nil {
		return []CommentDTO{}, fmt.Errorf("FindAllByNewsIDUseCase.Execute: %w", err)
	}

	tree := buildCommentTree(comments)

	return MapTreeToDTO(tree), nil
}

// buildCommentTree строит иерархию комментариев.
func buildCommentTree(comments []*dom.Comment) []*dom.Comment {
	idMap := make(map[int64]*dom.Comment)
	var roots []*dom.Comment

	// Индексируем
	for _, c := range comments {
		idMap[c.ID().Value()] = c
	}

	for _, c := range comments {
		parentID := c.ParentID().Value()
		if parentID != nil {
			if parent, ok := idMap[*parentID]; ok {
				parent.AddChild(c)
				continue
			}
		}
		roots = append(roots, c)
	}

	// Рекурсивно собираем детей
	for _, root := range roots {
		sortChildren(root)
	}

	return roots
}

// sortChildren сортирует вложенные комментарии.
func sortChildren(c *dom.Comment) {
	children := c.Children()

	sort.Slice(
		children, func(i, j int) bool {
			return children[i].PubTime().Time().Before(children[j].PubTime().Time())
		},
	)

	for _, child := range children {
		sortChildren(child)
	}
}
