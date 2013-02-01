package infrastructure

import (
	"code.google.com/p/sadbox/gettext"
	"os"
	"path"
	"path/filepath"
)

type GettextHandler struct {
	Catalogs           map[string]*gettext.Catalog
	supportedLanguages []string
}

func NewGettextHandler() (*GettextHandler, error) {
	gettextHandler := new(GettextHandler)
	gettextHandler.supportedLanguages = []string{"default", "english"}
	gettextHandler.Catalogs = make(map[string]*gettext.Catalog, len(gettextHandler.supportedLanguages))
	moReader := new(gettext.MoReader)
	for _, lang := range gettextHandler.supportedLanguages {
		b, err := os.Open(filepath.Join(".", filepath.FromSlash(path.Clean("/locale/"+lang+".mo"))))
		if err != nil {
			return nil, err
		}
		cat := gettext.NewCatalog()
		err = moReader.Read(cat, b)
		if err != nil {
			return nil, err
		}
		gettextHandler.Catalogs[lang] = cat
	}
	return gettextHandler, nil
}

func (handler *GettextHandler) Get(lang, key string, a ...interface{}) string {
	if cat, ok := handler.Catalogs[lang]; ok {
		return cat.Get(key, a...)
	}
	panic("Specified language not found")
	return ""
}

func (handler *GettextHandler) GetPlural(lang, key string, num int, a ...interface{}) string {
	if cat, ok := handler.Catalogs[lang]; ok {
		return cat.GetPlural(key, num, a...)
	}
	panic("Specified language not found")
	return ""
}

func (handler *GettextHandler) GetSupportedLanguages() []string {
	return handler.supportedLanguages
}
