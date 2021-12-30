package model

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

//UUID creates UUID for a session
func UUID() string {
	uuid := uuid.NewV4().String()
	uuid = strings.Replace(uuid, "-", "", -1)
	return uuid
}
