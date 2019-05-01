package notifier

import (
	"context"

	"github.com/858chain/token-shout/utils"
)

var EventTypeRegistry = []string{}

// Engine of notification
type Notifier struct {
	receivers map[string]*Receiver
	eventChan chan Event

	stopCh chan struct{}
}

// event generation buffer
const EventBufSize = 2 << 4

func New() *Notifier {
	return &Notifier{
		receivers: make(map[string]*Receiver),
		eventChan: make(chan Event, EventBufSize)}
}

func (notifier *Notifier) InstallReceiver(name string, receiver *Receiver) {
	notifier.receivers[name] = receiver
}

func (notifier *Notifier) UninstallReceiver(name string) {
	delete(notifier.receivers, name)
}

func (notifier *Notifier) ListReceivers() map[string]*Receiver {
	return notifier.receivers
}

func (notifier *Notifier) EventChan() chan<- Event {
	return notifier.eventChan
}

func (notifier *Notifier) Stop() {
	close(notifier.stopCh)
}

func (notifier *Notifier) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-notifier.stopCh:
			return

		case event, _ := <-notifier.eventChan:
			for name, receiver := range notifier.receivers {
				if receiver.Match(event) {
					utils.L.Debugf(name, "match", event)
					receiver.Accept(event)
				}
			}
		}
	}
}
