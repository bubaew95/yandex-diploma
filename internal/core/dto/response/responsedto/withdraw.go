package responsedto

type Withdraw struct {
	OrderNumber int64   `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
