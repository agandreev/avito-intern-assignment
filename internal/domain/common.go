package domain

const (
	AmountMode SortingMode = "amount"
	DateMode   SortingMode = "date"
)

// SortingMode represent string which describes Operation's order.
type SortingMode string

// ErrorJSON represents service error as struct for convenient response representation.
type ErrorJSON struct {
	Message string `json:"error"`
}

// OperationInput represents user's input for any operation except history.
type OperationInput struct {
	InitiatorID int64   `json:"initiator_id"`
	ReceiverID  int64   `json:"receiver_id"`
	Amount      float64 `json:"amount"`
}

// HistoryInput represents user's input for history operation.
type HistoryInput struct {
	ID       int64       `json:"id"`
	Quantity int64       `json:"quantity"`
	Mode     SortingMode `json:"mode"`
}
