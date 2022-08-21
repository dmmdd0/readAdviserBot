package telegram

import (
	"errors"
	"log"
	"net/url"
	"readAdviserBot/lib/e"
	"readAdviserBot/storage"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	//	add page: http://
	if isAddCmd(text) {
		//	TODO: AddPage()
	}

	//	rnd page: /rnd
	// help: /help
	//	start: /start: hi + help
	switch text {
	case RndCmd:
	case HelpCmd:
	case StartCmd:

	}
}

//todo: lesson #5 15:30
//	https://www.youtube.com/watch?v=f_esRaDae44&list=PLFAQFisfyqlWDwouVTUztKX2wUjYQ4T3l&index=5
func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreayExist)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	// todo: closure variant lesson #5 20:00
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p Processor) SendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't sens random", err) }()

	page, err := p.storage.PickRandom(username)
	//files not best plase for ErrNoSavePage because it is TG dipendent
	//if err != nil && !errors.Is(err, files.ErrNoSavedPage) {
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPage)
	}

}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	//	ya.ru don't valid for this type of cheking
	//	http(s)://ya.ru valid because prefix
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
