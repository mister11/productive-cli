package application

type LoginManager interface {
	IsSessionValid() (bool, error)
	Login(username string, password string) error
}