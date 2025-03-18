package systemdto

type CalculationSystem struct {
	Order   int64  `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
