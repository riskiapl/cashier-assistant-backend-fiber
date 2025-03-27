package types

import "time"

type MemberResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        *string   `json:"name"`
	PhoneNumber *string   `json:"phone_number"`
	Address     *string   `json:"address"`
}

type UpdateMemberRequest struct {
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Status      string  `json:"status"`
	Avatar      string  `json:"avatar"`
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
}

type MembersResponse struct {
	Members []MemberResponse `json:"members"`
	Count   int64            `json:"count"`
}
