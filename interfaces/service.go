package interfaces

import (
	"github.com/TheOnly92/morioka/usecases"
)

var gService *GlobalService

type GlobalService struct {
	UserInteractor *usecases.UserInteractor
	Translate      TranslateHandler
	Session        SessionHandler
}

func SetGlobalService(g *GlobalService) {
	gService = g
}
