package user

type BasicUser struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}
type Interface interface {
	CreateUser(BasicUser) (string, error)
}

type Service struct {
	repo
}
