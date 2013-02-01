package interfaces

import (
	"fmt"
	"net/http"
	"strings"
)

type Title struct {
	titles    []string
	separator string
	r         *http.Request
}

func NewTitle(separator string, r *http.Request) *Title {
	return &Title{separator: separator, r: r}
}

func (t *Title) Prepend(title string) {
	title = GetTranslate(t.r, title)
	t.titles = append([]string{title}, t.titles...)
}

func (t *Title) Append(title string) {
	title = GetTranslate(t.r, title)
	t.titles = append(t.titles, title)
}

func (t *Title) String() string {
	return strings.Join(t.titles, t.separator)
}

func GetLanguage(r *http.Request) string {
	cookie, err := r.Cookie("language")
	if err != nil {
		return "default"
	}
	switch cookie.Value {
	case "en":
		return "en"
	}
	return "default"
}

func GetTranslate(r *http.Request, key string, a ...interface{}) string {
	rt := gService.Translate.Get(GetLanguage(r), key, a...)
	if rt == "" {
		return fmt.Sprintf(key, a...)
	}
	return rt
}
