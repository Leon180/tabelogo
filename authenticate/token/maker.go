package token

import (
	"authenticate/model/entitymodel"
	"time"
)

type Maker interface {
	CreateToken(user entitymodel.User, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
