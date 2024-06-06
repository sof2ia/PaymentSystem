package internal

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) (int64, error)
}

type repository struct {
	client pgxClient
}

func (r *repository) CreateUser(ctx context.Context, user User) (int64, error) {
	row := QueryRow(ctx, r.client, `INSERT INTO Users ("name", "age", "phone", "email", "cpf", "balance") VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		user.Name, user.Age, user.Phone, user.Email, user.CPF, user.Balance)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int64(id), err
}

func NewRepository(client pgxClient) Repository {
	return &repository{client: client}
}
