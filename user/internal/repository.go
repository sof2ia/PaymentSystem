package internal

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) error
}

type repository struct {
	client pgxClient
}

func (r *repository) CreateUser(ctx context.Context, user User) error {
	_, err := Exec(ctx, `INSERT INTO Users ("name", "age", "phone", "email", "cpf") VALUES ($1, $2, $3, $4, $5)`,
		r.client, user.Name, user.Age, user.Phone, user.Email, user.CPF)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(client pgxClient) Repository {
	return &repository{client: client}
}
