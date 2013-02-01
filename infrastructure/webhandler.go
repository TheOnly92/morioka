package infrastructure

import (
	"encoding/json"
	"github.com/TheOnly92/morioka/interfaces"
	"github.com/TheOnly92/morioka/usecases"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
)

type JsonError struct {
	Error JsonErrorChild `json:"error"`
}

type JsonErrorChild struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *usecases.User, *interfaces.Title) *interfaces.WebError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				debug.PrintStack()
				page, err := template.New("error.html").ParseFiles("views/error.html")
				if err != nil {
					http.Error(w, "Internal Server Error", 500)
					panic(err)
				}
				err = page.Execute(w, map[string]interface{}{
					"Code":    500,
					"Message": "Internal Server Error!",
				})
				if err != nil {
					w.Write([]byte("Internal Server Error"))
					panic(err)
				}
				return
			}
		}()
		user, _ := interfaces.GetAuthUser(r)
		title := interfaces.NewTitle(" › ", r)
		title.Append("Project Morioka")
		if e := fn(w, r, user, title); e != nil {
			log.Println(e.Error)
			if e.Json {
				w.WriteHeader(e.Code)
				b, _ := json.Marshal(JsonError{JsonErrorChild{e.Code, e.Message}})
				w.Write(b)
			} else {
				page, err := template.New("error.html").ParseFiles("views/error.html")
				if err != nil {
					http.Error(w, "Internal Server Error", 500)
					panic(err)
				}
				err = page.Execute(w, map[string]interface{}{
					"Code":    e.Code,
					"Message": e.Message,
				})
				if err != nil {
					w.Write([]byte("Internal Server Error"))
					panic(err)
				}
			}
		}
	}
}
