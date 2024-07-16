package internal

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

type GetPixKeyResponse struct {
	UserID   int
	Name     string
	CPF      string
	KeyID    string
	KeyValue string
}
