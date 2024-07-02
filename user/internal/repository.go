package internal

import (
	"context"
	"log"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) (int64, error)
	GetUser(ctx context.Context, int642 int64) (user User, err error)
}

type repository struct {
	client pgxClient
}

func (r *repository) GetUser(ctx context.Context, idUser int64) (User, error) {
	row := QueryRow(ctx, r.client, `SELECT * FROM Users WHERE ID = $1`, idUser)
	log.Println(row)
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Email,
		&user.Phone,
		&user.CPF,
	)
	if err != nil {
		return User{}, err
	}
	log.Printf("id: %+v", user)
	return *user, nil
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
