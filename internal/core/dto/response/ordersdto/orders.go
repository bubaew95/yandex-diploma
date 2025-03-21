package ordersdto

type Orders struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}
