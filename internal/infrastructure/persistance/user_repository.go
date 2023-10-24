package persistance

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	"github.com/Deve-Lite/DashboardX-API/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/pkg/errors"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) bool {
	sqls := `
		SELECT "id" FROM "users" WHERE "email" = $1
	`

	var count int
	r.db.QueryRowxContext(ctx, sqls, email).Scan(&count)

	return count != 0
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	sqls := `
		SELECT "id", "name", "email", "password", "is_admin", "language", "theme"
		FROM "users" 
		WHERE "email" = $1
	`

	if err := r.db.GetContext(ctx, user, sqls, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ae.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "userRepository.GetByEmail.GetContext")
	}

	return user, nil
}

func (r *userRepository) Get(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user := &domain.User{}

	sqls := `
		SELECT "id", "name", "email", "password", "is_admin", "language", "theme"
		FROM "users" 
		WHERE "id" = $1
	`

	if err := r.db.GetContext(ctx, user, sqls, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ae.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "userRepository.Get.GetContext")
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.CreateUser) (uuid.UUID, error) {
	created := &domain.User{}
	f := `"name", "password", "email", "is_admin"`
	v := fmt.Sprintf(`'%s', '%s', '%s', %v`, user.Name, user.Password, user.Email, user.IsAdmin)

	if user.Language != "" {
		f = fmt.Sprintf(`%s, "language"`, f)
		v = fmt.Sprintf(`%s, '%s'`, v, user.Language)
	}

	if user.Theme != "" {
		f = fmt.Sprintf(`%s, "theme"`, f)
		v = fmt.Sprintf(`%s, '%s'::"theme"`, v, user.Theme)
	}

	sql := fmt.Sprintf(`
		INSERT INTO "users" (%s)
		VALUES (%s)
		RETURNING "id"
	`, f, v)

	if err := r.db.GetContext(ctx, created, sql); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == postgres.DuplicatedKey && pgErr.Constraint == postgres.UserEmailConstraint {
				return uuid.Nil, ae.ErrEmailExists
			}
		}

		return uuid.Nil, errors.Wrap(err, "userRepository.Create.GetContext")
	}

	return created.ID, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.UpdateUser) error {
	var p []string

	if user.Email.Set && !user.Email.Null {
		p = append(p, fmt.Sprintf(`"email" = '%s'`, user.Email.String))
	}

	if user.Name.Set && !user.Name.Null {
		p = append(p, fmt.Sprintf(`"name" = '%s'`, user.Name.String))
	}

	if user.IsAdmin.Set && !user.IsAdmin.Null {
		p = append(p, fmt.Sprintf(`"is_admin" = '%v'`, user.IsAdmin.Bool))
	}

	if user.Language.Set && !user.Language.Null {
		p = append(p, fmt.Sprintf(`"language" = '%s'`, user.Language.String))
	}

	if user.Theme.Set && !user.Theme.Null {
		p = append(p, fmt.Sprintf(`"theme" = '%s'`, user.Theme.String))
	}

	if user.Password.Set && !user.Password.Null {
		p = append(p, fmt.Sprintf(`"password" = '%s'`, user.Password.String))
	}

	if len(p) == 0 {
		return ae.ErrMissingParams
	}

	sql := fmt.Sprintf(`UPDATE "users" SET %s WHERE "id" = '%s'`, strings.Join(p, ","), user.ID.String())

	sr, err := r.db.ExecContext(ctx, sql)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == postgres.DuplicatedKey && pgErr.Constraint == postgres.UserEmailConstraint {
				return ae.ErrEmailExists
			}
		}

		return errors.Wrap(err, "userRepository.Update.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	sql := `DELETE FROM "users" WHERE "id" = $1`

	sr, err := r.db.ExecContext(ctx, sql, userID)
	if err != nil {
		return errors.Wrap(err, "userRepository.Delete.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrUserNotFound
	}

	return nil
}
