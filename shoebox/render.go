package shoebox

import (
	"errors"
	"os"
)

type ShoeboxRenderer struct {
}

func NewShoeboxRenderer() (*ShoeboxRenderer, error) {

	sb := ShoeboxRenderer{}

	return &sb, nil
}

func (sb *ShoeboxRenderer) RenderArchive(root_path string) error {

	info, err := os.Stat(root_path)

	if os.IsNotExist(err) {
		return err
	}

	if !info.IsDir() {
		return errors.New("Not a directory")
	}

	return nil
}
