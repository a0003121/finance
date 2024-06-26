package login

type Service interface {
	Login(username, password string) (string, error)
}
