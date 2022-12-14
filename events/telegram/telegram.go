package telegram

import (
	"errors"
	"readAdviserBot/clients/telegram"
	"readAdviserBot/events"
	"readAdviserBot/lib/e"
	"readAdviserBot/storage"
)

// not perfect name because interferate
type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrUnknounType = errors.New("unknown event type")
var ErrUnknounMetaType = errors.New("uhknown meta typw")

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

//func (p Processor) Fetch(limit int) ([]events.Event, error) {
//	update, err := p.tg.Updates(p.offset, limit)
//	if err != nil {
//		return nil, e.Wrap("can't get events", err)
//	}
//
//	res := make([]events.Event, 0, len(update))
//
//	for _, u := range update {
//		res = append(res, event(u))
//	}
//}

// last updated Fetch
func (p Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get updates", err)
	}
	res := make([]events.Event, 0, len(updates))

	if len(updates) == 0 {
		return nil, nil
	}

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMesage(event)
	default:
		return e.Wrap("can't proccess message", ErrUnknounType)
	}
}

func (p *Processor) processMesage(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		//todo: remove repit message
		return e.Wrap("can't process message", err)

	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknounMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	// chatID and username telegram only parameters
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Usernamr,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Uncnown
	}

	return events.Message

}
