package utility

import "github.com/google/uuid"

func GenDefaultUUID() string {
	return uuid.NewString()
}
