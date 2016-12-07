package tool

import (
	"github.com/satori/go.uuid"
)

func GenUUID() string {
	taksId := uuid.NewV4()
	return taksId.String()
}
