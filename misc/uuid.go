package misc

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func UUID() uuid.UUID {
	return uuid.NewV4()
}

func UUIDString() string {
	return fmt.Sprintf("%v", UUID())
}
