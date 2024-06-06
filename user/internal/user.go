package internal

import (
	pb "PaymentSystem/protobuf"
)

type User struct {
	ID      int64
	Name    string
	Age     int32
	Phone   string
	Email   string
	CPF     string
	Balance float64
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

type CreateUserResponse struct {
	ID int64
}
