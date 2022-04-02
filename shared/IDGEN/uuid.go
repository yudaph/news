package IDGEN

import "github.com/google/uuid"

var NewUUID = func() string {
	return uuid.NewString()
}
