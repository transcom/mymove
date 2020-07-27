package models

import (
	"encoding/json"
)

// ServerNotification contains the message sent to the server
type ServerNotification struct {
	Message string
}

// String is not required by pop and may be deleted
func (e ServerNotification) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}
