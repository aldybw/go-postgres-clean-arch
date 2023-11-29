package usecase

import (
	"context"
	"go-postgres-clean-arch/domain"
	"time"

	"github.com/go-playground/validator"
)

type tagUsecase struct {
	tagRepo        domain.TagRepository
	contextTimeout time.Duration
	validate       *validator.Validate
}

func NewTagUsecase(t domain.TagRepository, timeout time.Duration, v *validator.Validate) domain.TagUseCase {
	return &tagUsecase{
		tagRepo:        t,
		contextTimeout: timeout,
		validate:       v,
	}
}

// Fetch implements domain.TagUseCase.
func (t *tagUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Tag, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	res, nextCursor, err = t.tagRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

// FetchByID implements domain.TagUseCase.
func (t *tagUsecase) FetchByID(c context.Context, id int64) (res domain.Tag, err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	res, err = t.tagRepo.FetchByID(ctx, id)
	if err != nil {
		return
	}

	return
}

// FetchByName implements domain.TagUseCase.
func (t *tagUsecase) FetchByName(c context.Context, name string) (res domain.Tag, err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	res, err = t.tagRepo.FetchByName(ctx, name)
	if err != nil {
		return
	}

	return
}

// Store implements domain.TagUseCase.
func (t *tagUsecase) Store(c context.Context, tag *domain.Tag) (err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	existedTag, _ := t.FetchByName(ctx, tag.Name)
	// check conflict tag
	if existedTag != (domain.Tag{}) {
		return domain.ErrConflict
	}

	err = t.tagRepo.Store(ctx, tag)
	return
}

// Update implements domain.TagUseCase.
func (t *tagUsecase) Update(c context.Context, tag *domain.Tag) (err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	tag.UpdatedAt = time.Now()
	return t.tagRepo.Update(ctx, tag)
}

// Delete implements domain.TagUseCase.
func (t *tagUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (domain.Article{}) {
		return domain.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}

// OLD IMPLEMENTATION
// // Delete implements domain.TagUseCase.
// func (t *tagUsecase) Delete(tagId int) error {
// 	err := t.tagRepository.Delete(tagId)

// 	if err != nil {
// 		return errors.New(err.Error())
// 	}
// 	return nil
// }

// // Fetch implements domain.TagUseCase.
// func (t *tagUsecase) Fetch() ([]domain.TagResponse, error) {
// 	result, err := t.tagRepository.Fetch()

// 	var tags []domain.TagResponse
// 	for _, value := range result {
// 		tag := &domain.TagResponse{
// 			ID:   value.ID,
// 			Name: value.Name,
// 		}
// 		tags = append(tags, *tag)
// 	}

// 	if err != nil {
// 		return nil, errors.New(err.Error())
// 	}
// 	return tags, nil
// }

// // FetchByID implements domain.TagUseCase.
// func (t *tagUsecase) FetchByID(tagId int) (domain.TagResponse, error) {
// 	tagData, err := t.tagRepository.FetchByID(tagId)

// 	tagResponse := &domain.TagResponse{
// 		ID:   tagData.ID,
// 		Name: tagData.Name,
// 	}

// 	if err != nil {
// 		return *tagResponse, errors.New(err.Error())
// 	}

// 	return *tagResponse, nil
// }

// // Store implements domain.TagUseCase.
// func (t *tagUsecase) Store(tag *domain.CreateTagRequest) error {
// 	err := t.validate.Struct(tag)
// 	if err != nil {
// 		return errors.New(err.Error())
// 	}

// 	tagDomain := domain.Tag{
// 		Name: tag.Name,
// 	}

// 	errRepo := t.tagRepository.Store(&tagDomain)
// 	if errRepo != nil {
// 		return errors.New(errRepo.Error())
// 	}
// 	return nil
// }

// // Update implements domain.TagUseCase.
// func (t *tagUsecase) Update(tag *domain.UpdateTagRequest) error {
// 	tagData, err := t.tagRepository.FetchByID(tag.ID)
// 	if err != nil {
// 		return errors.New(err.Error())
// 	}

// 	tagData.Name = tag.Name
// 	t.tagRepository.Update(&tagData)

// 	return nil
// }
