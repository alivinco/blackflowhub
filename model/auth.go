package model

type AuthRequest struct {
	IsAuthenticated bool
	Username string
	Email string
	Error error
}
