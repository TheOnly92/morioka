package domain

import ()

type Account struct {
	Id             int
	Type           int
	Name           string
	Order          int
	StartingAmount int64
	Comment        string
	OwnerId        int

	TypeName   string
	CreditCard *CreditCard
}

type AccountType struct {
	Id   int
	Name string
}

type CreditCard struct {
	LastDate      int
	PayingMonth   int
	PayingDay     int
	PayingAccount int
	Holiday       int
}

type AccountRepository interface {
	FindById(id, userId int) (*Account, error)
	FindAllByOwner(userId int) ([]*Account, error)
	FindAllTypes() ([]*AccountType, error)
	FindAllByType(accType, userId int) ([]*Account, error)
	Store(account *Account) error
	Delete(account *Account) error
}
