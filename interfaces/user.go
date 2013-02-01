package interfaces

import (
	"github.com/TheOnly92/morioka/usecases"
	"github.com/gorilla/sessions"
	"net/http"
)

func GetAuthUser(r *http.Request) (*usecases.User, bool) {
	id, log := GetAuthId(r)
	if !log {
		user := &usecases.User{}
		return user, false
	}
	user, err := gService.UserInteractor.FindById(id)
	if err != nil {
		user = &usecases.User{}
		return user, false
	}
	return user, true
}

func GetAuthId(r *http.Request) (int, bool) {
	sessionAuth, _ := GetSessionStore().Get(r, "auth")
	var authIdT interface{}
	var ok bool
	if authIdT, ok = sessionAuth.Values["id"]; !ok {
		return 0, false
	}
	return authIdT.(int), true
}

func GetSessionStore() sessions.Store {
	return gService.Session.GetSessionStore()
}
