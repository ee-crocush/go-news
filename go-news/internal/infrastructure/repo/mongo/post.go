package mongo

import (
	dom "GoNews/internal/domain/post"
	"GoNews/internal/infrastructure/repo/mongo/mapper"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

var _ dom.Repository = (*PostRepository)(nil)

// PostRepository представляет собой репозиторий для работы с новостями в MongoDB.
type PostRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
	timeout    time.Duration
}

// NewPostRepository создаёт новый Mongo-репозиторий с новостями.
func NewPostRepository(db *mongo.Database, timeout time.Duration) *PostRepository {
	return &PostRepository{
		db:         db,
		collection: db.Collection("posts"),
		timeout:    timeout,
	}
}

// Store сохраняет новости в MongoDB.
func (r *PostRepository) Store(ctx context.Context, post *dom.Post) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	id, err := r.getNextID(ctx)
	if err != nil {
		return fmt.Errorf("PostRepository.Create: %w", err)
	}

	postID, err := dom.NewPostID(id)
	if err != nil {
		return fmt.Errorf("PostRepository.Create: %w", err)
	}

	post.SetID(postID)

	doc := mapper.FromPostToDoc(post)

	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("PostRepository.Create: %w", err)
	}

	return nil
}

// FindByID находит новость по его ID.
func (r *PostRepository) FindByID(ctx context.Context, postID dom.PostID) (*dom.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var doc mapper.PostDocument

	if err := r.collection.FindOne(ctx, bson.M{"_id": postID.Value()}).Decode(&doc); err != nil {
		return nil, fmt.Errorf("PostRepository.FindByID: %w", err)
	}

	return mapper.MapDocToPost(doc)
}

// FindLast получает последнюю новость.
func (r *PostRepository) FindLast(ctx context.Context) (*dom.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var doc mapper.PostDocument

	opts := options.FindOne().SetSort(bson.D{{Key: "pub_time", Value: -1}})
	if err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&doc); err != nil {
		return nil, fmt.Errorf("PostRepository.FindLast: %w", err)
	}

	return mapper.MapDocToPost(doc)
}

// FindLatest получает последние n новостей.
func (r *PostRepository) FindLatest(ctx context.Context, limit int) ([]*dom.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "pub_time", Value: -1}}).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("PostRepository.FindLatest: %w", err)
	}
	defer cursor.Close(ctx)

	return r.decodeManyPosts(ctx, cursor)
}

// FindAll получает все новости.
func (r *PostRepository) FindAll(ctx context.Context) ([]*dom.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("PostRepository.FindAll: %w", err)
	}
	defer cursor.Close(ctx)

	return r.decodeManyPosts(ctx, cursor)
}

// getNextID возвращает следующее значение идентификатора.
func (r *PostRepository) getNextID(ctx context.Context) (int32, error) {
	filter := bson.M{"_id": "posts"}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		Seq int32 `bson:"seq"`
	}
	err := r.db.Collection("counters").FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("getNextID: %w", err)
	}

	return result.Seq, nil
}

// decodeManyPosts декодирует курсор в массив новостей.
func (r *PostRepository) decodeManyPosts(ctx context.Context, cursor *mongo.Cursor) ([]*dom.Post, error) {
	var posts []*dom.Post

	for cursor.Next(ctx) {
		var doc mapper.PostDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("PostRepository.decodeManyPosts: %w", err)
		}

		post, err := mapper.MapDocToPost(doc)
		if err != nil {
			return nil, fmt.Errorf("PostRepository.decodeManyPosts: %w", err)
		}

		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("PostRepository.decodeManyPosts: %w", err)
	}

	return posts, nil
}
