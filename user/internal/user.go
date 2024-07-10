package internal

import (
	pb "PaymentSystem/protobuf"
)

type User struct {
	ID         int
	Name       string
	Age        int32
	Phone      string
	Email      string
	CPF        string
	Balance    float64
	ListPixKey []PixKey
}

type CreateUserRequest struct {
	Name  string `validate:"required"`
	Age   int32  `validate:"required,gte=18"`
	Phone string `validate:"required,phone"`
	Email string `validate:"required,email"`
	CPF   string `validate:"required,CPF"`
}

func ConvertCreateUserRequest(userPB *pb.CreateUserRequest) (CreateUserRequest, error) {
	c := CreateUserRequest{
		Name:  userPB.Name,
		Age:   userPB.Age,
		Phone: userPB.Phone,
		Email: userPB.Email,
		CPF:   userPB.Cpf,
	}
	err := c.Validation()
	if err != nil {
		return CreateUserRequest{}, err
	}
	return c, nil
}

func ConvertGetUserResponse(user *User) (*pb.GetUserResponse, error) {
	g := &pb.GetUserResponse{
		Name:    user.Name,
		Age:     user.Age,
		Phone:   user.Phone,
		Email:   user.Email,
		Cpf:     user.CPF,
		Balance: user.Balance,
	}
	return g, nil
}

type PixKey struct {
	KeyID    string
	UserID   int
	KeyType  KeyType
	KeyValue string
}

type KeyType string

const (
	Phone  KeyType = "phone"
	Email  KeyType = "email"
	CPF    KeyType = "cpf"
	Random KeyType = "random"
)
