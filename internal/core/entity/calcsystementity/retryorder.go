package calcsystementity

import (
	"github.com/bubaew95/yandex-diploma/internal/core/entity/orderentity"
	"time"
)

type Worker struct {
	RetryQueueCh chan orderentity.OrderDetails
	RetryTimer   time.Time
	Order        orderentity.OrderDetails
}
