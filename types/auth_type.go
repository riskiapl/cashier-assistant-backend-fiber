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

type RegisterInput struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PlainPassword string `json:"plain_password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}
