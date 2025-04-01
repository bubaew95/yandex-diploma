package calcsystementity

import (
	"github.com/bubaew95/yandex-diploma/internal/core/entity/orderentity"
	"time"
)

type RetryOrder struct {
	Order     orderentity.OrderDetails
	RetryTime time.Time
}
