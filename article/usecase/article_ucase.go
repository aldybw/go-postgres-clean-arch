package usecase

import (
	"context"
	"go-postgres-clean-arch/domain"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type articleUsecase struct {
	articleRepo    domain.ArticleRepository
	tagRepo        domain.TagRepository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewArticleUsecase(a domain.ArticleRepository, t domain.TagRepository, timeout time.Duration) domain.ArticleUsecase {
	return &articleUsecase{
		articleRepo:    a,
		tagRepo:        t,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *articleUsecase) fillTagDetails(c context.Context, data []domain.Article) ([]domain.Article, error) {
	g, ctx := errgroup.WithContext(c)

	// Get the tag's id
	mapTags := map[int64]domain.Tag{}

	for _, article := range data { //nolint
		mapTags[article.Tag.ID] = domain.Tag{}
	}
	// Using goroutine to fetch the tag's detail
	chanTag := make(chan domain.Tag)
	for tagID := range mapTags {
		tagID := tagID
		g.Go(func() error {
			res, err := a.tagRepo.FetchByID(ctx, tagID)
			if err != nil {
				return err
			}
			chanTag <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanTag)
	}()

	for tag := range chanTag {
		if tag != (domain.Tag{}) {
			mapTags[tag.ID] = tag
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the tag's data
	for index, item := range data { //nolint
		if a, ok := mapTags[item.Tag.ID]; ok {
			data[index].Tag = a
		}
	}
	return data, nil
}

func (a *articleUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillTagDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *articleUsecase) GetByID(c context.Context, id int64) (res domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resTag, err := a.tagRepo.FetchByID(ctx, res.Tag.ID)
	if err != nil {
		return domain.Article{}, err
	}
	res.Tag = resTag
	return
}

func (a *articleUsecase) Update(c context.Context, ar *domain.UpdateArticleInput) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	selectedArticle, err := a.GetByID(ctx, ar.ID)
	if err != nil {
		return err
	}

	existedArticle, _ := a.GetByTitle(ctx, ar.Title)

	if existedArticle.Title != selectedArticle.Title && existedArticle.Title == ar.Title {
		return domain.ErrConflict
	}

	if ar.Title == "" {
		ar.Title = selectedArticle.Title
	}
	if ar.Content == "" {
		ar.Content = selectedArticle.Content
	}
	if ar.TagID == 0 {
		ar.TagID = selectedArticle.Tag.ID
	}

	_, err = a.tagRepo.FetchByID(ctx, ar.TagID)
	if err != nil {
		return err
	}

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleUsecase) GetByTitle(c context.Context, title string) (res domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resTag, err := a.tagRepo.FetchByID(ctx, res.Tag.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Tag = resTag
	return
}

func (a *articleUsecase) Store(c context.Context, m *domain.CreateArticleInput) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, _ := a.GetByTitle(ctx, m.Title)

	if existedArticle.Title == m.Title {
		return domain.ErrConflict
	}

	_, err = a.tagRepo.FetchByID(ctx, m.TagID)
	if err != nil {
		return err
	}

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *articleUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
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
