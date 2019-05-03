package notifier

import (
	"encoding/json"
)

type Event interface {
	GetEvent() map[string]interface{}

	// eth_balance_change_event / erc20_token_transfer
	Type() string
	From() string
	To() string
}

type Msg struct {
	EventType string                 `json:"type"`
	Event     map[string]interface{} `json:"event"`
}

func EncodeEvent(event Event) ([]byte, error) {
	var msg = &Msg{
		EventType: event.Type(),
		Event:     event.GetEvent(),
	}

	return json.Marshal(msg)
}
