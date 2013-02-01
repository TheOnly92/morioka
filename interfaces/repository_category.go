package interfaces

import (
	"github.com/TheOnly92/morioka/domain"
)

type DbCategoryRepo DbRepo

func NewDbCategoryRepo(db DbHandler) *DbCategoryRepo {
	return &DbCategoryRepo{
		db: db,
	}
}

func (repo *DbCategoryRepo) FindAllCategories(catType string, userId int) ([]*domain.Category, error) {
	rows, err := repo.db.Query("SELECT id, name, category_type, owner_id, display_order FROM categories WHERE category_type = $1 AND owner_id = $2 ORDER BY display_order ASC", catType, userId)
	if err != nil {
		return nil, err
	}
	var rt []*domain.Category
	for rows.Next() {
		category := new(domain.Category)
		err = rows.Scan(&category.Id, &category.Name, &category.Type, &category.OwnerId, &category.Order)
		if err != nil {
			return nil, err
		}
		rt = append(rt, category)
	}
	return rt, nil
}

func (repo *DbCategoryRepo) FindById(id, userId int) (*domain.Category, error) {
	rt := new(domain.Category)
	err := repo.db.QueryRow("SELECT id, name, category_type, owner_id, display_order FROM categories WHERE id = $1 AND owner_id = $2", id, userId).Scan(&rt.Id, &rt.Name, &rt.Type, &rt.OwnerId, &rt.Order)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (repo *DbCategoryRepo) FindDetailById(detailId, catId, userId int) (*domain.CategoryDetail, error) {
	rt := new(domain.CategoryDetail)
	err := repo.db.QueryRow("SELECT id, category_id, name, owner_id, display_order FROM category_details WHERE category_id = $1 AND id = $2 AND owner_id = $3", catId, detailId, userId).Scan(&rt.Id, &rt.CategoryId, &rt.Name, &rt.OwnerId, &rt.Order)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (repo *DbCategoryRepo) FindAllDetailsByCategory(catId, userId int) ([]*domain.CategoryDetail, error) {
	rows, err := repo.db.Query("SELECT id, category_id, name, owner_id, display_order FROM category_details WHERE category_id = $1 AND owner_id = $2 ORDER BY category_id ASC, display_order ASC", catId, userId)
	if err != nil {
		return nil, err
	}
	var rt []*domain.CategoryDetail
	for rows.Next() {
		detail := new(domain.CategoryDetail)
		err = rows.Scan(&detail.Id, &detail.CategoryId, &detail.Name, &detail.OwnerId, &detail.Order)
		if err != nil {
			return nil, err
		}
		rt = append(rt, detail)
	}
	return rt, nil
}

func (repo *DbCategoryRepo) FindAllDetails(userId int) ([]*domain.CategoryDetail, error) {
	rows, err := repo.db.Query("SELECT id, category_id, name, owner_id, display_order FROM category_details WHERE owner_id = $1 ORDER BY category_id ASC, display_order ASC", userId)
	if err != nil {
		return nil, err
	}
	var rt []*domain.CategoryDetail
	for rows.Next() {
		detail := new(domain.CategoryDetail)
		err = rows.Scan(&detail.Id, &detail.CategoryId, &detail.Name, &detail.OwnerId, &detail.Order)
		if err != nil {
			return nil, err
		}
		rt = append(rt, detail)
	}
	return rt, nil
}

func (repo *DbCategoryRepo) Store(category *domain.Category) error {
	if category.Id == 0 {
		err := repo.db.QueryRow("INSERT INTO categories (name, category_type, owner_id, display_order) VALUES ($1, $2, $3, (SELECT COALESCE(MAX(display_order), 0) FROM categories WHERE category_type = $2 AND owner_id = $3)) RETURNING id", category.Name, category.Type, category.OwnerId).Scan(&category.Id)
		if err != nil {
			return err
		}
	} else {
		_, err := repo.db.Exec("UPDATE categories SET name = $1, category_type = $2, display_order = $3 WHERE id = $4 AND owner_id = $5", category.Name, category.Type, category.Order, category.Id, category.OwnerId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *DbCategoryRepo) StoreDetail(detail *domain.CategoryDetail) error {
	if detail.Id == 0 {
		err := repo.db.QueryRow("INSERT INTO category_details (category_id, name, display_order, owner_id) VALUES ($1, $2, (SELECT MAX(display_order) FROM category_details WHERE category_id = $1 AND owner_id = $3), $3) RETURNING id", detail.CategoryId, detail.Name, detail.Order, detail.OwnerId).Scan(&detail.Id)
		if err != nil {
			return err
		}
	} else {
		_, err := repo.db.Exec("UPDATE category_details SET name = $1, category_id = $2, display_order = $3 WHERE id = $4 AND owner_id = $5", detail.Name, detail.CategoryId, detail.Order, detail.Id, detail.OwnerId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *DbCategoryRepo) Delete(category *domain.Category) error {
	_, err := repo.db.Exec("DELETE FROM categories WHERE id = $1 AND owner_id = $2", category.Id, category.OwnerId)
	if err != nil {
		return err
	}
	return nil
}
