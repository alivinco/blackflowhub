package model

type UserAuth struct {
	ID     string `json:"id"`
	Login  string `json:"login"`
	Token  string `json:"-"`
	Secret string `json:"-"`
	Expiry int64  `json:"-"`
	Email  string `json:"email"`
	Avatar string `json:"avatar_url"`
	Active bool   `json:"active,"`
	Admin  bool   `json:"admin,"`
	Hash   string `json:"-"`
}