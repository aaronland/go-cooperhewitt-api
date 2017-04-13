package shoebox

import (
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-cooperhewitt-api/template"
	"github.com/thisisaaronland/go-cooperhewitt-api/util"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
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

	index_items := make([]*template.ShoeboxIndexItem, 0)

	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	callback := func(abs_path string, info os.FileInfo) error {

		if info.IsDir() {
			return nil
		}

		fname := filepath.Base(abs_path)

		if fname != "index.json" {
			return nil
		}

		// log.Println("archive", abs_path)

		wg.Add(1)
		defer wg.Done()

		fh, err := os.Open(abs_path)

		if err != nil {
			return err
		}

		defer fh.Close()

		item, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		i, err := sb.RenderItem(root_path, item)

		if err != nil {
			log.Println(fmt.Sprintf("failed to render %s, because %s", abs_path, err))
			return nil

			// return err
		}

		mu.Lock()
		index_items = append(index_items, i)
		mu.Unlock()

		return nil
	}

	c := crawl.NewCrawler(root_path)

	err = c.Crawl(callback)

	if err != nil {
		return err
	}

	wg.Wait()

	err = sb.RenderIndex(root_path, index_items)

	if err != nil {
		return err
	}

	return nil
}

func (sb *ShoeboxRenderer) RenderIndex(root_path string, sb_items []*template.ShoeboxIndexItem) error {

	fname := "index.html"
	abs_path := filepath.Join(root_path, fname)

	fh, err := os.Create(abs_path)

	if err != nil {
		return err
	}

	defer fh.Close()

	n := "index"
	t, err := template.NewShoeboxIndex(n)

	if err != nil {
		return err
	}

	sb_data := template.ShoeboxIndex{
		Items: sb_items,
	}

	err = t.ExecuteTemplate(fh, n, sb_data)

	if err != nil {
		return err
	}

	return nil
}

func (sb *ShoeboxRenderer) RenderItem(root_path string, item []byte) (*template.ShoeboxIndexItem, error) {

	item_id := gjson.GetBytes(item, "id").Int()
	item_title := gjson.GetBytes(item, "title").String()

	path := util.Id2Path(item_id)
	rel_path := filepath.Join(root_path, path)

	refersto_id := gjson.GetBytes(item, "refers_to_uid").Int()

	refersto_fname := fmt.Sprintf("%d.json", refersto_id)
	refersto_path := filepath.Join(rel_path, refersto_fname)

	// TO DO: handle things that are not objects...

	fh, err := os.Open(refersto_path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	object, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	fname := "index.html"
	abs_path := filepath.Join(rel_path, fname)

	t_fh, err := os.Create(abs_path)

	if err != nil {
		return nil, err
	}

	defer t_fh.Close()

	n := "item"
	t, err := template.NewShoeboxItem(n)

	if err != nil {
		return nil, err
	}

	obj_title := gjson.GetBytes(object, "object.title").String()
	obj_url := gjson.GetBytes(object, "object.url").String()

	sb_object := template.ShoeboxObject{
		Title: obj_title,
		URL:   obj_url,
	}

	sb_data := template.ShoeboxItem{
		Title:  obj_title,
		Object: sb_object,
	}

	err = t.ExecuteTemplate(t_fh, n, sb_data)

	if err != nil {
		return nil, err
	}

	item_url := filepath.Join(path, fname)

	i := template.ShoeboxIndexItem{
		Id:    item_id,
		URL:   item_url,
		Title: item_title,
	}

	return &i, nil
}
