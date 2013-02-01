package usecases

import (
	"github.com/TheOnly92/morioka/domain"
)

type AccountInteractor struct {
	AccountRepo domain.AccountRepository
}

func (inter *AccountInteractor) List(userId int) ([]*domain.Account, error) {
	return inter.AccountRepo.FindAllByOwner(userId)
}

func (inter *AccountInteractor) ListTypes() ([]*domain.AccountType, error) {
	return inter.AccountRepo.FindAllTypes()
}

func (inter *AccountInteractor) FetchById(id, userId int) (*domain.Account, error) {
	return inter.AccountRepo.FindById(id, userId)
}

func (inter *AccountInteractor) Save(account *domain.Account) error {
	return inter.AccountRepo.Store(account)
}

func (inter *AccountInteractor) Delete(account *domain.Account) error {
	return inter.AccountRepo.Delete(account)
}
