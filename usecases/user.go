package usecases

type User struct {
	Id       int
	Email    string
	Password string
}

type UserRepository interface {
	FindById(id int) (*User, error)
}

type UserInteractor struct {
	UserRepo UserRepository
}

func (inter *UserInteractor) FindById(id int) (*User, error) {
	return inter.UserRepo.FindById(id)
}
