package interfaces

import (
	"github.com/TheOnly92/morioka/usecases"
)

type DbUserRepo DbRepo

func NewDbUserRepo(db DbHandler) *DbUserRepo {
	return &DbUserRepo{
		db: db,
	}
}

func (repo *DbUserRepo) FindById(id int) (*usecases.User, error) {
	rt := new(usecases.User)
	err := repo.db.QueryRow("SELECT id, email, password FROM users WHERE id = $1", id).Scan(&rt.Id, &rt.Email, &rt.Password)
	if err != nil {
		return nil, err
	}
	return rt, nil
}
