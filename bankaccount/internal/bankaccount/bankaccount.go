package bankaccount

type TransferRequest struct {
	PayerID        string  `validate:"required"`
	ReceiverPixKey string  `validate:"required"`
	Amount         float64 `validate:"required,gt=0,lte=5000"`
}

type TransferRespond struct {
	Success bool
}

type Movement struct {
	ID            string
	Amount        string
	UserID        string
	Date          string
	TransactionID string
	OperationType string
}
