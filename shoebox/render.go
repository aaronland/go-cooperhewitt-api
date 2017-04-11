package shoebox

import (
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-cooperhewitt-api/util"
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

		fname := filepath.Base(abs_path)

		if fname != "info.json" {
			return nil
		}

		fh, err := os.Open(abs_path)

		if err != nil {
			return err
		}

		defer fh.Close()

		item, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		err = sb.RenderItem(root_path, item)

		if err != nil {
			return err
		}

		return nil
	}

	c := crawl.NewCrawler(root_path)

	err = c.Crawl(callback)

	if err != nil {
		return err
	}

	return nil
}

func (sb *ShoeboxRenderer) RenderItem(root_path string, item []byte) error {

	item_id := gjson.GetBytes(item, "id").Int()

	path := util.Id2Path(item_id)
	rel_path := filepath.Join(root_path, path)

	refersto_id := gjson.GetBytes(item, "refers_to_uid").Int()

	refersto_fname := fmt.Sprintf("%d.json", refersto_id)
	refersto_path := filepath.Join(rel_path, refersto_fname)

	log.Println(refersto_path)
	return nil
}
