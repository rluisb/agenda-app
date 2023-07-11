package protocols

type EmailValidator interface {
	IsValid(email string) bool
}