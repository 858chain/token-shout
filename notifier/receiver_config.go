package notifier

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/858chain/token-shout/utils"
)

// Config of receivers
type ReceiverConfig struct {
	RetryCount          uint     `json:"retrycount"`
	Endpoint            string   `json:"endpoint"`
	NewBalanceRemaining float64  `json:"newBalanceRemaining"`
	EventTypes          []string `json:"eventTypes"`

	// address support the following format
	//
	// * - all valid address
	// all-in-wallet - all address in local wallet
	// 0x9d8ec6a2cfe2a6abdf5a6041e752ba851abb4dbb,0x9d8ec6a2cfe2a6abdf5a6041e752ba851abb4dbb
	FromAddresses []string `json:"from"`
	ToAddresses   []string `json:"to"`
}

func (rc ReceiverConfig) ValidCheck() error {
	if rc.RetryCount <= 0 {
		return errors.New("retryCount should greater than 0")
	}

	if _, err := url.Parse(rc.Endpoint); err != nil {
		return err
	}

	if (rc.NewBalanceRemaining) < 0 {
		return errors.New("newBalanceRemaining should greater than 0")
	}

	if len(rc.EventTypes) == 0 {
		return errors.New("eventTypes not provided")
	}

	// EventType should registered.
	for _, etype := range rc.EventTypes {
		if !utils.StringSliceContains(EventTypeRegistry, etype) {
			return errors.New(fmt.Sprintf("%s not a valid event type, make sure in %+v", etype, EventTypeRegistry))
		}
	}

	return nil
}
