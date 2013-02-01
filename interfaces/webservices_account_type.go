package interfaces

import (
	"github.com/TheOnly92/morioka/usecases"
	"encoding/json"
	"net/http"
)

type AccountTypeHandler struct {
	AccountInteractor *usecases.AccountInteractor
}

func (handler AccountTypeHandler) GetRoutes() []*RouteHandler {
	return []*RouteHandler{
		{"/account-types/list.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.ListJSON(w, r, user, title)
		}, []string{"GET"}},
	}
}

func (handler AccountTypeHandler) ListJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	types, _ := handler.AccountInteractor.ListTypes()
	b, _ := json.Marshal(types)
	w.Write(b)
	return nil
}
