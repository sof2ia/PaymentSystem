package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) (int64, error)
}

type repository struct {
	db *pgx.Conn
}

func (r *repository) CreateUser(ctx context.Context, user User) (int64, error) {
	resp := r.db.QueryRow(ctx, `INSERT INTO Users ("name", "age", "phone", "email", "cpf") VALUES (?, ?, ?, ?, ?) RETURNING id`,
		user.Name, user.Age, user.Phone, user.Email, user.CPF)
	var id int64
	err := resp.Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func NewRepository(db *pgx.Conn) Repository {
	return &repository{db: db}
}
