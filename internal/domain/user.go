package domain

import (
	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Password string    `db:"password"`
	Email    string    `db:"email"`
	IsAdmin  bool      `db:"is_admin"`
	Language string    `db:"language"`
	Theme    string    `db:"theme"`
}

type CreateUser struct {
	Name     string  `db:"name"`
	Password string  `db:"password"`
	Email    string  `db:"email"`
	IsAdmin  bool    `db:"is_admin"`
	Language *string `db:"language"`
	Theme    *string `db:"theme"`
}

type UpdateUser struct {
	ID       uuid.UUID `db:"id"`
	Name     t.String  `db:"name"`
	Password t.String  `db:"password"`
	Email    t.String  `db:"email"`
	IsAdmin  t.Bool    `db:"is_admin"`
	Language t.String  `db:"language"`
	Theme    t.String  `db:"theme"`
}
