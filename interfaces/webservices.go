package interfaces

import (
	"github.com/TheOnly92/morioka/domain"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
)

type Layout struct {
	Menus         []*MenuItem
	Title         string
	BootstrapData map[string]interface{}
	BuildVersion  string
	Router        string
}

type RenderMenu struct {
	Menus           []*MenuItem
	DisplayLanguage string
}

func BackboneTemplate(bootstrap map[string]interface{}, router string, w http.ResponseWriter, title *Title, identifier string) error {
	t, err := template.New("backbone-layout.html").ParseFiles("views/backbone-layout.html")
	if err != nil {
		return err
	}
	err = t.Execute(w, &Layout{
		ConstructMenu(identifier),
		title.String(),
		bootstrap,
		domain.BuildVersion,
		router,
	})
	if err != nil {
		return err
	}
	return nil
}

func ParseJSONBody(r *http.Request) (map[string]interface{}, error) {
	if r.Method != "POST" && r.Method != "PUT" {
		return nil, errors.New("method not supported")
	}
	var rt map[string]interface{}
	var err error
	ct := r.Header.Get("Content-Type")
	ct, _, err = mime.ParseMediaType(ct)
	if ct == "application/json" {
		var reader io.Reader = r.Body
		maxFormSize := int64(1<<63 - 1)
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		if int64(len(b)) > maxFormSize {
			return nil, errors.New("POST/PUT too large")
		}
		err = json.Unmarshal(b, &rt)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func InterfaceToInt(input interface{}, def int) int {
	if s, ok := input.(float64); ok {
		return int(s)
	} else if s, ok := input.(string); ok {
		rt, _ := strconv.Atoi(s)
		return rt
	} else if s, ok := input.(int); ok {
		return s
	} else if s, ok := input.(int64); ok {
		return int(s)
	}
	return def
}

func InterfaceToInt64(input interface{}, def int64) int64 {
	if s, ok := input.(float64); ok {
		return int64(s)
	} else if s, ok := input.(string); ok {
		rt, _ := strconv.ParseInt(s, 10, 64)
		return rt
	} else if s, ok := input.(int); ok {
		return int64(s)
	} else if s, ok := input.(int64); ok {
		return s
	}
	return def
}
