package interfaces

import (
	"encoding/json"
	"github.com/TheOnly92/morioka/domain"
	"github.com/TheOnly92/morioka/usecases"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	CategoryInteractor *usecases.CategoryInteractor
}

func (handler CategoryHandler) GetRoutes() []*RouteHandler {
	return []*RouteHandler{
		{"/categories/", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.Manage(w, r, user, title)
		}, []string{"GET"}},
		{"/categories/list.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.ListJSON(w, r, user, title)
		}, []string{"GET"}},
		{"/categories/{id:[0-9]+}.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.FetchJSON(w, r, user, title)
		}, []string{"GET"}},
		{"/categories/{id:[0-9]+|create}.json", func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
			return handler.PutPostJSON(w, r, user, title)
		}, []string{"PUT", "POST"}},
	}
}

func (handler CategoryHandler) Manage(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	if user.Id == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	title.Prepend("項目登録")
	categories, err := handler.CategoryInteractor.ListAll(user.Id)
	if err != nil {
		return &WebError{
			Error:   err,
			Code:    500,
			Message: "Error retrieving categories",
		}
	}
	err = BackboneTemplate(map[string]interface{}{
		"Categories": categories,
	}, "routers/category", w, title, "CategoryManageHandler")
	if err != nil {
		return &WebError{
			Error:   err,
			Code:    500,
			Message: "Error rendering template",
		}
	}
	return nil
}

func (handler CategoryHandler) FetchJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	category, err := handler.CategoryInteractor.FetchById(id, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}
	b, _ := json.Marshal(category)
	w.Write(b)
	return nil
}

func (handler CategoryHandler) ListJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	categories, _ := handler.CategoryInteractor.ListAll(user.Id)
	b, _ := json.Marshal(categories)
	w.Write(b)
	return nil
}

func (handler CategoryHandler) PutPostJSON(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError {
	vars := mux.Vars(r)
	var err error
	category := new(domain.Category)
	if vars["id"] != "create" {
		id, _ := strconv.Atoi(vars["id"])
		category, err = handler.CategoryInteractor.FetchById(id, user.Id)
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
		category.Name = name
	}
	category.Order = InterfaceToInt(rt["Order"], category.Order)
	err = handler.CategoryInteractor.Save(category)
	if err != nil {
		return &WebError{
			Error:   err,
			Code:    500,
			Message: "Error saving category",
			Json:    true,
		}
	}
	return nil
}
