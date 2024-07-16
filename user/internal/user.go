package internal

import (
	pb "PaymentSystem/protobuf"
	"strconv"
)

type User struct {
	ID      int
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

func ConvertUpdateUserRequest(user *pb.UpdateUserRequest) (User, error) {
	idUserInt, err := strconv.Atoi(user.UserId)
	if err != nil {
		return User{}, err
	}
	resPB := User{
		ID:    idUserInt,
		Name:  user.Name,
		Age:   user.Age,
		Phone: user.Phone,
		Email: user.Email,
		CPF:   user.Cpf,
	}
	return resPB, nil
}
func ConvertUpdateUserResponse(user *User) (*pb.UpdateUserResponse, error) {
	resPB := &pb.UpdateUserResponse{
		UserId: strconv.Itoa(user.ID),
		Name:   user.Name,
		Age:    user.Age,
		Phone:  user.Phone,
		Email:  user.Email,
		Cpf:    user.CPF,
	}
	return resPB, nil
}
