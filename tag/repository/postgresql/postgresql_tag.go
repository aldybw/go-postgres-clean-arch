package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"go-postgres-clean-arch/domain"
	"go-postgres-clean-arch/tag/repository"
	"time"

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
	var query string
	if cursor != "" {
		query = `SELECT id,name,created_at,updated_at
					FROM tag
					WHERE created_at > $1
					ORDER BY created_at
					LIMIT $2`

		decodedCursor, err := repository.DecodeCursor(cursor)
		if err != nil && cursor != "" {
			return nil, "", domain.ErrBadParamInput
		}
		res, err = p.fetch(ctx, query, decodedCursor, num)
		if err != nil {
			return nil, "", err
		}

		// fmt.Println("len(res): " + strconv.Itoa(len(res)))
		// fmt.Println("int(num): " + strconv.Itoa(int(num)))
		if len(res) == int(num) {
			nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
			// decodedCursor, err := repository.DecodeCursor(cursor)
			// decodedCursor, err = repository.DecodeCursor(nextCursor)
			if err != nil && cursor != "" {
				return nil, "", domain.ErrBadParamInput
			}
			// fmt.Println("res[len(res)-1].CreatedAt: " + res[len(res)-1].CreatedAt.String())
			// fmt.Println("cursor: " + cursor)
			// fmt.Println("EncodeCursor: " + nextCursor)
			// fmt.Println("decodedCursor: " + decodedCursor.String())
		}

		return res, nextCursor, nil
	}

	query = `SELECT id,name,created_at,updated_at
			FROM tag
			ORDER BY created_at
			LIMIT $1`
	res, err = p.fetch(ctx, query, num)
	if err != nil {
		return nil, "", err
	}

	return
}

// FetchByID implements domain.TagRepository.
func (p *postgresqlTagRepo) FetchByID(ctx context.Context, id int64) (res domain.Tag, err error) {
	query := `SELECT id,name,created_at,updated_at 
				FROM tag 
				WHERE id = $1`

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
	query := `SELECT id, name, created_at, updated_at 
				FROM tag 
				WHERE name = $1`

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
	query := `INSERT INTO tag (name, created_at, updated_at) 
				VALUES ($1, $2, $3)
				RETURNING id`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.QueryRowContext(ctx, t.Name, time.Now(), time.Now()).Err()
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	return
}

// Delete implements domain.TagRepository.
func (p *postgresqlTagRepo) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM tag WHERE id = $1"

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
func (p *postgresqlTagRepo) Update(ctx context.Context, id int64, t *domain.Tag) (err error) {
	query := `UPDATE tag SET name=$1, updated_at=$2 WHERE id = $3;`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, t.Name, t.UpdatedAt, id)
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
