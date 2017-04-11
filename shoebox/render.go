package shoebox

import (
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ShoeboxRenderItem struct {
	Id    int64
	Title string
}

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

	// item_ids := []*ShoeboxRenderItem

	callback := func(abs_path string, info os.FileInfo) error {

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(abs_path)

		if ext != ".json" {
			return nil
		}

		fh, err := os.Open(abs_path)

		if err != nil {
			return err
		}

		defer fh.Close()

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		log.Println(gjson.GetBytes(body, "id"))
		return nil
	}

	c := crawl.NewCrawler(root_path)

	err = c.Crawl(callback)

	if err != nil {
		return err
	}

	return nil
}
