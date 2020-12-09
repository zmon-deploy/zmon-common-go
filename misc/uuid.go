package misc

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func UUID() uuid.UUID {
	uid, _ := uuid.NewV4()
	return uid
}

func UUIDString() string {
	return fmt.Sprintf("%v", UUID())
}
