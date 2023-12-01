package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"go-postgres-clean-arch/article/repository"
	"go-postgres-clean-arch/domain"

	"github.com/sirupsen/logrus"
)

type postgresqlArticleRepository struct {
	Conn *sql.DB
}

// NewPostgresqlArticleRepository will create an object that represent the article.Repository interface
func NewPostgresqlArticleRepository(conn *sql.DB) domain.ArticleRepository {
	return &postgresqlArticleRepository{conn}
}

func (m *postgresqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Article, 0)
	for rows.Next() {
		t := domain.Article{}
		tagID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&tagID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Tag = domain.Tag{
			ID: tagID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *postgresqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	var query string

	if cursor != "" {
		query = `SELECT id,title,content, tag_id, updated_at, created_at
					FROM article 
					WHERE created_at > $1 
					ORDER BY created_at 
					LIMIT $2 `

		// Example decodedCursor or cursor: 2023-11-30 15:20:33.682 +0000 UTC
		// Example encodeCursor and cursor: MjAyMy0xMS0zMFQxNToyMDozMy42ODJa

		decodedCursor, err := repository.DecodeCursor(cursor)
		if err != nil && cursor != "" {
			return nil, "", domain.ErrBadParamInput
		}

		res, err = m.fetch(ctx, query, decodedCursor, num)
		if err != nil {
			return nil, "", err
		}

		if len(res) == int(num) {
			nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
		}

		return res, nextCursor, nil
	}

	query = `SELECT id,title,content, tag_id, updated_at, created_at
				FROM article 
				ORDER BY created_at 
				LIMIT $1 `

	res, err = m.fetch(ctx, query, num)
	if err != nil {
		return nil, "", err
	}

	return
}

func (m *postgresqlArticleRepository) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	query := `SELECT id,title,content, tag_id, updated_at, created_at
				FROM article 
				WHERE ID = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *postgresqlArticleRepository) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	query := `SELECT id,title,content, tag_id, updated_at, created_at
				FROM article 
				WHERE title = $1`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *postgresqlArticleRepository) Store(ctx context.Context, a *domain.ArticleInput) (err error) {
	query := `INSERT INTO article (title, content, tag_id, updated_at , created_at) 
				VALUES ($1, $2, $3, $4, $5)
				RETURNING ID`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.QueryRowContext(ctx, a.Title, a.Content, a.TagID, a.UpdatedAt, a.CreatedAt).Err()
	if err != nil {
		return
	}
	return
}

func (m *postgresqlArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM article WHERE id = $1"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}

func (m *postgresqlArticleRepository) Update(ctx context.Context, ar *domain.ArticleInput) (err error) {
	query := `UPDATE article SET title=$1, content=$2, tag_id=$3, updated_at=$4 WHERE id = $5;`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.TagID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
