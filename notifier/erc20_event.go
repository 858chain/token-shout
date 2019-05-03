package notifier

const (
	EventNameERC20 = "erc20_log_event"
)

func init() {
	EventTypeRegistry = append(EventTypeRegistry, EventNameERC20)
}

type ERC20LogEvent struct {
	meta map[string]interface{}
}

func NewERC20LogEvent(meta map[string]interface{}) *ERC20LogEvent {
	return &ERC20LogEvent{
		meta: meta,
	}
}

func (ethEvent *ERC20LogEvent) GetEvent() map[string]interface{} {
	return ethEvent.meta
}

func (ethEvent *ERC20LogEvent) Type() string {
	return EventNameERC20
}

func (ethEvent *ERC20LogEvent) From() string {
	if from, found := ethEvent.meta["from"]; !found {
		return ""
	} else {
		return from.(string)
	}
}

func (ethEvent *ERC20LogEvent) To() string {
	if to, found := ethEvent.meta["to"]; !found {
		return ""
	} else {
		return to.(string)
	}
}
