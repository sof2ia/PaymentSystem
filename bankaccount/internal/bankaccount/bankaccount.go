package bankaccount

type TransferRequest struct {
	FromUserID string
	ToPixKey   string
	Amount     float64
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
