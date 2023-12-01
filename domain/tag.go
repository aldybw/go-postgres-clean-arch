package domain

import (
	"context"
	"time"
)

// Tag is representing the Tag data struct
type Tag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TagUseCase represent the tag's usecases
type TagUseCase interface {
	Fetch(ctx context.Context, cursor string, num int64) (tags []Tag, nextCursor string, err error) // naked return
	FetchByID(ctx context.Context, id int64) (Tag, error)
	FetchByName(ctx context.Context, name string) (Tag, error)
	Store(ctx context.Context, t *Tag) error
	Update(ctx context.Context, id int64, t *Tag) error
	Delete(ctx context.Context, id int64) error
}

// Tag represent the tag's repository contract
type TagRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (tags []Tag, nextCursor string, err error) // naked return
	FetchByID(ctx context.Context, id int64) (Tag, error)
	FetchByName(ctx context.Context, name string) (Tag, error)
	Store(ctx context.Context, t *Tag) error
	Update(ctx context.Context, id int64, t *Tag) error
	Delete(ctx context.Context, id int64) error
}

// // CreateTagRequest is representing the create request data input
// type CreateTagRequest struct {
// 	Name string `validate:"required,min=1,max=200" json:"name"`
// }

// // UpdateTagRequest is representing the update request data input
// type UpdateTagRequest struct {
// 	ID   int    `validate:"required"`
// 	Name string `validate:"required,max=200,min=1" json:"name"`
// }

// // TagResponse is representing how the response should be presented
// type TagResponse struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

