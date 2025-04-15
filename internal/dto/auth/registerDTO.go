package dto

// Пока что это существует просто для адекватной работы сваггера, Некит - тебе надо бы с этим разобраться
type RegisterDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
