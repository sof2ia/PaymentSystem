package internal

import (
	"context"
	"log"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) (int, error)
	GetUser(ctx context.Context, idUser int) (user *User, err error)
	CreatePixKey(ctx context.Context, pix PixKey) (string, error)
}

type repository struct {
	client pgxClient
}

func (r *repository) GetUser(ctx context.Context, idUser int) (*User, error) {
	row := QueryRow(ctx, r.client, `SELECT * FROM Users WHERE ID = $1`, idUser)
	log.Printf("row: %+v", row)
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
		return nil, err
	}
	log.Printf("user: %+v", *user)
	return user, nil
}

func (r *repository) CreateUser(ctx context.Context, user User) (int, error) {
	row := QueryRow(ctx, r.client, `INSERT INTO Users ("name", "age", "phone", "email", "cpf", "balance") VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		user.Name, user.Age, user.Phone, user.Email, user.CPF, user.Balance)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (r *repository) CreatePixKey(ctx context.Context, pix PixKey) (string, error) {
	row := QueryRow(ctx, r.client, `INSERT INTO PixKey("user_id", "key_type", "key_value") VALUES (?, ?, ?)`,
		pix.UserID, pix.KeyType, pix.KeyValue)
	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func NewRepository(client pgxClient) Repository {
	return &repository{client: client}
}
