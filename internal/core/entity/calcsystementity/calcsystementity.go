package calcsystementity

import "github.com/bubaew95/yandex-diploma/internal/core/dto/response/systemdto"

type CalculationSystem struct {
	*systemdto.CalculationSystem
	StatusCode int
	Retry      string
	UserId     int64
}
