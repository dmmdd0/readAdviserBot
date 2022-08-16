package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"readAdviserBot/lib/e"
	"readAdviserBot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defoultPerm = 0774

var ErNoSavedPage = errors.New("no saved page")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defoultPerm); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	// ignore error if bad Close
	defer func() { _ = file.Close() }()

	//	serialice page
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErNoSavedPage
	}

	//	random from 0 to quantity of files
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.PickRandom()
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return e.Wrap(msg, err)
	}
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exist", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	// lesson #4 bug report :=
	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exist", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

//	open and decode
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()
	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	return &p, nil
}

// function need for future, when you will change hash type
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
