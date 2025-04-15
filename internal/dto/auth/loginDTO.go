package dto

// Пока что это существует просто для адекватной работы сваггера, Некит - тебе надо бы с этим разобраться
type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
