package interfaces

type Service interface {
	Register()
	Login()
}

type Repository interface {
	CreateUser()
	GetUserByName()
}
