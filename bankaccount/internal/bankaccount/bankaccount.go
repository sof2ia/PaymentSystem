package bankaccount

import pb "PaymentSystem/protobuf"

type TransferRequest struct {
	PayerID        string  `validate:"required"`
	ReceiverPixKey string  `validate:"required"`
	Amount         float64 `validate:"required,gt=0,lte=5000"`
}

type TransferResponse struct {
	Success bool
}

type Movement struct {
	ID            string
	Amount        string
	UserID        string
	Date          string
	TransactionID string
	OperationType OperationType
}

type OperationType string

const (
	Credit OperationType = "credit"
	Debit  OperationType = "debit"
)

func ConvertTransferRequest(requestPB *pb.TransferRequest) TransferRequest {
	return TransferRequest{
		PayerID:        requestPB.FromUserId,
		ReceiverPixKey: requestPB.ToPixKey,
		Amount:         requestPB.Amount,
	}
}
