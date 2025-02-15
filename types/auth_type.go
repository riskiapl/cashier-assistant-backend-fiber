package types

type LoginInput struct {
	Userormail string `json:"userormail"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
