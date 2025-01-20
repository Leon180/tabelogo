package inject

import (
	"authenticate/token"

	"github.com/google/wire"
)

var tokenMakerHandleSet = wire.NewSet(
	provideTokenMaker,
)

func provideTokenMaker(
	symmetricKey string,
) token.Maker {
	tokenMaker, _ := token.NewPasetoMaker(symmetricKey)
	return tokenMaker
}
