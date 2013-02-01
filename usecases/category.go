package usecases

import (
	"github.com/TheOnly92/morioka/domain"
)

type CategoryInteractor struct {
	CategoryRepo domain.CategoryRepository
}

func (inter *CategoryInteractor) ListAll(userId int) ([]*domain.Category, error) {
	rt, err := inter.CategoryRepo.FindAllCategories("expense", userId)
	if err != nil {
		return nil, err
	}
	t, err := inter.CategoryRepo.FindAllCategories("income", userId)
	if err != nil {
		return nil, err
	}
	rt = append(rt, t...)
	return rt, nil
}

func (inter *CategoryInteractor) FetchById(id, userId int) (*domain.Category, error) {
	return inter.CategoryRepo.FindById(id, userId)
}

func (inter *CategoryInteractor) Save(category *domain.Category) error {
	return inter.CategoryRepo.Store(category)
}
