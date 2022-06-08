package tinkoff

type PaymentStatus string

// статусы платежа:
// https://www.tinkoff.ru/kassa/develop/api/payments/
const (
	StatusNew             PaymentStatus = "NEW"              // Платеж создан
	StatusDeadlineExpired PaymentStatus = "DEADLINE_EXPIRED" // Время оплаты истекло
	StatusAttemptsExpired PaymentStatus = "ATTEMPTS_EXPIRED" // Попытки открытия формы исчерпаны
	StatusFormShowed      PaymentStatus = "FORM_SHOWED"      // Платежная форма открыта покупателем
	StatusAuthorized      PaymentStatus = "AUTHORIZED"       // Денежный средства зарезервированы
	StatusRejected        PaymentStatus = "REJECTED"         // Платеж отменен банком
	StatusConfirmed       PaymentStatus = "CONFIRMED"        // Подтвержден
	StatusRefunded        PaymentStatus = "REFUNDED"         // Денежные средства возвращены полностью
)
