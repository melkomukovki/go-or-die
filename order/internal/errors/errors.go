package errs

import "errors"

var (
	ErrInvalidUUID                 = errors.New("неверный формат UUID")
	ErrOrderNotFound               = errors.New("заказ не найден")
	ErrOrderAlreadyPaid            = errors.New("заказ уже оплачен")
	ErrOrderCancelled              = errors.New("заказ отменён")
	ErrPartNotFound                = errors.New("деталь не найдена")
	ErrOutOfStock                  = errors.New("деталь отсутствует на складе")
	ErrOrderPendingPaymentMismatch = errors.New("заказ не в статусе ожидания оплаты")
)
