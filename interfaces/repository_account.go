package interfaces

import (
	"github.com/TheOnly92/morioka/domain"
)

type DbAccountRepo DbRepo

func NewDbAccountRepo(db DbHandler) *DbAccountRepo {
	return &DbAccountRepo{
		db: db,
	}
}

func (repo *DbAccountRepo) FindById(id, userId int) (*domain.Account, error) {
	rt := new(domain.Account)
	card := new(domain.CreditCard)
	err := repo.db.QueryRow("SELECT a.id, account_type, account_name, starting_amount, at.name, display_order, COALESCE(comment, ''), COALESCE(c.last_date, 0), COALESCE(c.paying_month, 0), COALESCE(c.paying_day, 0), COALESCE(c.paying_account, 0), COALESCE(c.holiday, 0) FROM accounts AS a INNER JOIN account_types AS at ON (at.id = a.account_type) LEFT OUTER JOIN credit_cards AS c ON (c.account_id = a.id) WHERE a.id = $1 AND a.owner_id = $2", id, userId).Scan(&rt.Id, &rt.Type, &rt.Name, &rt.StartingAmount, &rt.TypeName, &rt.Order, &rt.Comment, &card.LastDate, &card.PayingMonth, &card.PayingDay, &card.PayingAccount, &card.Holiday)
	if err != nil {
		return nil, err
	}
	if rt.Type == 4 {
		rt.CreditCard = card
	}
	return rt, nil
}

func (repo *DbAccountRepo) FindAllByOwner(userId int) ([]*domain.Account, error) {
	rows, err := repo.db.Query("SELECT a.id, account_type, account_name, starting_amount, at.name, display_order, COALESCE(c.last_date, 0), COALESCE(c.paying_month, 0), COALESCE(c.paying_day, 0), COALESCE(c.paying_account, 0), COALESCE(c.holiday, 0) FROM accounts AS a INNER JOIN account_types AS at ON (at.id = a.account_type) LEFT OUTER JOIN credit_cards AS c ON (c.account_id = a.id) WHERE a.owner_id = $1 ORDER BY display_order ASC, id ASC", userId)
	if err != nil {
		return nil, err
	}

	var rt []*domain.Account
	for rows.Next() {
		account := new(domain.Account)
		card := new(domain.CreditCard)
		err = rows.Scan(&account.Id, &account.Type, &account.Name, &account.StartingAmount, &account.TypeName, &account.Order, &card.LastDate, &card.PayingMonth, &card.PayingDay, &card.PayingAccount, &card.Holiday)
		if err != nil {
			return nil, err
		}
		if account.Type == 4 {
			account.CreditCard = card
		}
		rt = append(rt, account)
	}
	return rt, nil
}

func (repo *DbAccountRepo) FindAllTypes() ([]*domain.AccountType, error) {
	rows, err := repo.db.Query("SELECT id, name FROM account_types ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	var rt []*domain.AccountType
	for rows.Next() {
		accountType := new(domain.AccountType)
		err = rows.Scan(&accountType.Id, &accountType.Name)
		if err != nil {
			return nil, err
		}
		rt = append(rt, accountType)
	}
	return rt, nil
}

func (repo *DbAccountRepo) FindAllByType(accType, userId int) ([]*domain.Account, error) {
	rows, err := repo.db.Query("SELECT a.id, account_type, account_name, starting_amount, display_order, at.name FROM accounts AS a INNER JOIN account_types AS at ON (at.id = a.account_type) WHERE account_type = $1 AND owner_id = $2 ORDER BY display_order ASC, id ASC", accType, userId)
	if err != nil {
		return nil, err
	}

	var rt []*domain.Account
	for rows.Next() {
		account := new(domain.Account)
		err = rows.Scan(&account.Id, &account.Type, &account.Name, &account.StartingAmount, &account.Order, &account.TypeName)
		if err != nil {
			return nil, err
		}
		rt = append(rt, account)
	}
	return rt, nil
}

func (repo *DbAccountRepo) Store(account *domain.Account) error {
	if account.Id == 0 {
		err := repo.db.QueryRow("INSERT INTO accounts (account_type, account_name, display_order, starting_amount, comment, owner_id) VALUES ($1, $2, (SELECT COALESCE(MAX(display_order), 0) + 1 FROM accounts WHERE owner_id = $5), $3, $4, $5) RETURNING id", account.Type, account.Name, account.StartingAmount, account.Comment, account.OwnerId).Scan(&account.Id)
		if err != nil {
			return err
		}
		if account.CreditCard != nil {
			_, err = repo.db.Exec("INSERT INTO credit_cards (account_id, last_date, paying_month, paying_day, paying_account, holiday, owner_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", account.Id, account.CreditCard.LastDate, account.CreditCard.PayingMonth, account.CreditCard.PayingDay, account.CreditCard.PayingAccount, account.CreditCard.Holiday, account.OwnerId)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := repo.db.Exec("UPDATE accounts SET account_type = $1, account_name = $2, starting_amount = $3, comment = $4, display_order = $7 WHERE id = $5 AND owner_id = $6", account.Type, account.Name, account.StartingAmount, account.Comment, account.Id, account.OwnerId, account.Order)
		if err != nil {
			return err
		}
		if account.CreditCard != nil {
			_, err = repo.db.Exec("UPDATE credit_cards SET last_date = $1, paying_month = $2, paying_day = $3, paying_account = $4, holiday = $5 WHERE account_id = $6 AND owner_id = $7", account.CreditCard.LastDate, account.CreditCard.PayingMonth, account.CreditCard.PayingDay, account.CreditCard.PayingAccount, account.CreditCard.Holiday, account.Id, account.OwnerId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *DbAccountRepo) Delete(account *domain.Account) error {
	_, err := repo.db.Exec("DELETE FROM accounts WHERE id = $1 AND owner_id = $2", account.Id, account.OwnerId)
	if err != nil {
		return err
	}
	return nil
}
