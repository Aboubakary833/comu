package domain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mazen160/go-random"
)

var (
	ErrPostNotFound    = errors.New("the post you're looking for does'nt exist")
	ErrCommentNotFound = errors.New("the comment you're looking for does'nt exist")
	ErrUnauthorized    = errors.New("sorry, you can't perform this action")
)

const DefaultPaginatorLimit = 10

type Post struct {
	ID        uuid.UUID	`json:"id"`
	UserID    uuid.UUID	`json:"user_id"`
	Title     string	`json:"title"`
	Slug      string	`json:"slug"`
	Content   string	`json:"content"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
}

type Comment struct {
	ID        uuid.UUID	`json:"id"`
	PostID    uuid.UUID	`json:"post_id"`
	UserID    uuid.UUID	`json:"user_id"`
	Content   string	`json:"content"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
}

type Cursor struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func (c *Cursor) ToBase64() (string, error) {
	rawBytes, err := json.Marshal(*c)

	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(rawBytes), nil
}

type Paginator struct {
	Limit int
	After *Cursor
}

func NewPost(authorID uuid.UUID, title, content string) *Post {
	slug := MakePostSlug(title)

	return &Post{
		UserID:    authorID,
		Title:     title,
		Slug:      slug,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewComment(postID, authorID uuid.UUID, content string) *Comment {
	return &Comment{
		PostID:    postID,
		UserID:    authorID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func MakePostSlug(title string) string {
	fields := strings.Fields(strings.ToLower(title))
	rand, _ := random.String(16)

	return strings.Join(fields, "-") + rand
}

type PostRepository interface {
	FindByID(context.Context, uuid.UUID) (*Post, error)
	FindBySlug(context.Context, string) (*Post, error)
	ListAll(context.Context) ([]Post, error)
	List(context.Context, Paginator) ([]Post, *Cursor, error)
	Store(context.Context, *Post) error
	Update(context.Context, *Post) error
	Delete(context.Context, *Post) error
}

type CommentRepository interface {
	Find(context.Context, uuid.UUID) (*Comment, error)
	ListAll(context.Context, uuid.UUID) ([]Comment, error)
	List(context.Context, uuid.UUID, Paginator) ([]Comment, *Cursor, error)
	Store(context.Context, *Comment) error
	Update(context.Context, *Comment) error
	Delete(context.Context, *Comment) error
}
