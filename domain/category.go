package domain

type Category struct {
	Id      int
	Name    string
	Type    string
	OwnerId int
	Order   int
}

type CategoryDetail struct {
	Id         int
	CategoryId int
	Name       string
	OwnerId    int
	Order      int
}

type CategoryRepository interface {
	FindAllCategories(catType string, userId int) ([]*Category, error)
	FindById(id, userId int) (*Category, error)
	FindDetailById(detailId, catId, userId int) (*CategoryDetail, error)
	FindAllDetailsByCategory(catId, userId int) ([]*CategoryDetail, error)
	FindAllDetails(userId int) ([]*CategoryDetail, error)
	Store(category *Category) error
	StoreDetail(detail *CategoryDetail) error
	Delete(category *Category) error
}
