package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"go-postgres-clean-arch/domain"
	"go-postgres-clean-arch/tag/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postgresqlTagRepo struct {
	Conn *sql.DB
	Db   *gorm.DB
}

// PostgresqlTagRepository will create an object that represent the tag.Repository interface
func NewPostgresqlTagRepository(Conn *sql.DB, Db *gorm.DB) domain.TagRepository {
	return &postgresqlTagRepo{Conn: Conn, Db: Db}
}

func (p *postgresqlTagRepo) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Tag, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
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

	result = make([]domain.Tag, 0)
	for rows.Next() {
		t := domain.Tag{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

// Fetch implements domain.TagRepository.
func (p *postgresqlTagRepo) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Tag, nextCursor string, err error) {
	query := `SELECT id,name,created_at,updated_at 
				FROM tags 
				WHERE created_at > ? 
				ORDER BY created_at 
				LIMIT ?`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = p.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

// FetchByID implements domain.TagRepository.
func (p *postgresqlTagRepo) FetchByID(ctx context.Context, id int64) (res domain.Tag, err error) {
	query := `SELECT id,name,created_at,updated_at 
				FROM tags 
				WHERE id = ?`

	list, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Tag{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

// FetchByName implements domain.TagRepository.
func (p *postgresqlTagRepo) FetchByName(ctx context.Context, name string) (res domain.Tag, err error) {
	query := `SELECT id,name,created_at,updated_at 
				FROM tags 
				WHERE name = ?`

	list, err := p.fetch(ctx, query, name)
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

// Store implements domain.TagRepository.
func (p *postgresqlTagRepo) Store(ctx context.Context, t *domain.Tag) (err error) {
	query := `INSERT  tags SET name=? , created_at=? , updated_at=?`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, t.Name, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	t.ID = lastID
	return
}

// Delete implements domain.TagRepository.
func (p *postgresqlTagRepo) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM tags WHERE id = ?"

	stmt, err := p.Conn.PrepareContext(ctx, query)
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

// Update implements domain.TagRepository.
func (p *postgresqlTagRepo) Update(ctx context.Context, t *domain.Tag) (err error) {
	query := `UPDATE tags set name=?, updated_at=? WHERE id = ?`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, t.Name, t.UpdatedAt, t.ID)
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

// OLD IMPLEMENTATION
// // Delete implements domain.TagRepository.
// func (p *postgresqlTagRepository) Delete(tagId int) error {
// 	var tag domain.Tag
// 	result := p.Db.Where("id = ?", tagId).Delete(&tag)
// 	// helper.ErrorPanic(result.Error)

// 	if result.Error != nil {
// 		return errors.New("failed to delete tag")
// 	}
// 	return nil

// }

// // Fetch implements domain.TagRepository.
// func (p *postgresqlTagRepository) Fetch() ([]domain.Tag, error) {
// 	var tags []domain.Tag
// 	result := p.Db.Find(&tags)
// 	// helper.ErrorPanic(result.Error)

// 	if result.Error != nil {
// 		return tags, errors.New("failed to get tags")
// 	}
// 	return tags, nil

// }

// // FetchByID implements domain.TagRepository.
// func (p *postgresqlTagRepository) FetchByID(tagId int) (domain.Tag, error) {
// 	var tag domain.Tag
// 	result := p.Db.Find(&tag, tagId)

// 	if result.Error != nil {
// 		return tag, errors.New("tag not found")
// 	}
// 	return tag, nil
// }

// // Store implements domain.TagRepository.
// func (p *postgresqlTagRepository) Store(t *domain.Tag) error {
// 	result := p.Db.Create(t)

// 	if result.Error != nil {
// 		return errors.New("failed to create tag")
// 	}
// 	return nil
// }

// // Update implements domain.TagRepository.
// func (p *postgresqlTagRepository) Update(t *domain.Tag) error {
// 	var updateTag = domain.UpdateTagRequest{
// 		ID:   t.ID,
// 		Name: t.Name,
// 	}
// 	result := p.Db.Model(&t).Updates(updateTag)

// 	if result.Error != nil {
// 		return errors.New("failed to update tag")
// 	}
// 	return nil
// }
