package interfaces

import (
	"github.com/TheOnly92/morioka/domain"
	"github.com/TheOnly92/morioka/usecases"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	AccountInteractor *usecases.AccountInteractor
}

func (handler AccountHandler) GetRoutes() []*RouteHandler {
	return []*RouteHandler{
		{"/accounts/", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.Manage(w, r, user, title)
		}, []string{"GET"}},
		{"/accounts/list.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.ListJSON(w, r, user, title)
		}, []string{"GET"}},
		{"/accounts/{id:[0-9]+}.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.FetchJSON(w, r, user, title)
		}, []string{"GET"}},
		{"/accounts/{id:[0-9]+|create}.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.PutPostJSON(w, r, user, title)
		}, []string{"PUT", "POST"}},
		{"/accounts/{id:[0-9]+}.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.DeleteJSON(w, r, user, title)
		}, []string{"DELETE"}},
	}
}

func (handler AccountHandler) Manage(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	if user.Id == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	title.Prepend("口座一覧")
	accounts, err := handler.AccountInteractor.List(user.Id)
	if err != nil {
		return &WebError{
			Error:   err,
			Code:    500,
			Message: "Error retrieving accounts",
		}
	}
	err = BackboneTemplate(map[string]interface{}{
		"Accounts": accounts,
	}, "routers/account", w, title, "AccountManageHandler")
	if err != nil {
		return &WebError{
			Error:   err,
			Code:    500,
			Message: "Error rendering template",
		}
	}
	return nil
}

func (handler AccountHandler) ListJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	accounts, _ := handler.AccountInteractor.List(user.Id)
	b, _ := json.Marshal(accounts)
	w.Write(b)
	return nil
}

func (handler AccountHandler) DeleteJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	account, err := handler.AccountInteractor.FetchById(id, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}
	account.OwnerId = user.Id
	err = handler.AccountInteractor.Delete(account)
	if err != nil {
		log.Println(err)
		return nil
	}
	return nil
}

func (handler AccountHandler) FetchJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	account, err := handler.AccountInteractor.FetchById(id, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}
	b, _ := json.Marshal(account)
	w.Write(b)
	return nil
}

func (handler AccountHandler) PutPostJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	vars := mux.Vars(r)
	var err error
	account := new(domain.Account)
	if vars["id"] != "create" {
		id, _ := strconv.Atoi(vars["id"])
		account, err = handler.AccountInteractor.FetchById(id, user.Id)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	rt, err := ParseJSONBody(r)
	if err != nil {
		log.Println(err)
		return nil
	}
	if name, ok := rt["Name"].(string); ok {
		account.Name = name
	}
	account.Type = InterfaceToInt(rt["Type"], 0)
	account.StartingAmount = InterfaceToInt64(rt["StartingAmount"], 0)
	account.Order = InterfaceToInt(rt["Order"], account.Order)
	if comment, ok := rt["Comment"].(string); ok {
		account.Comment = comment
	}
	if rt["CreditCard"] != nil {
		if card, ok := rt["CreditCard"].(map[string]interface{}); ok {
			if account.CreditCard == nil {
				account.CreditCard = new(domain.CreditCard)
			}
			account.CreditCard.LastDate = InterfaceToInt(card["LastDate"], 0)
			account.CreditCard.PayingMonth = InterfaceToInt(card["PayingMonth"], 0)
			account.CreditCard.PayingDay = InterfaceToInt(card["PayingDay"], 0)
			account.CreditCard.PayingAccount = InterfaceToInt(card["PayingAccount"], 0)
			account.CreditCard.Holiday = InterfaceToInt(card["Holiday"], 0)
		}
	}
	account.OwnerId = user.Id
	err = handler.AccountInteractor.Save(account)
	if err != nil {
		return &WebError{
			Error:   err,
			Code:    500,
			Message: "Error saving account",
			Json:    true,
		}
	}
	return nil
}
