// universal package for all messager
package events

// move offser onto Fetcher
type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Uncnown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
