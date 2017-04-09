package shoebox

import (
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-cooperhewitt-api"
	"github.com/thisisaaronland/go-cooperhewitt-api/util"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Shoebox struct {
	client api.APIClient
}

func NewShoebox(client api.APIClient) (*Shoebox, error) {

	sb := Shoebox{
		client: client,
	}

	return &sb, nil
}

func (sb *Shoebox) Archive(dest string) error {

	info, err := os.Stat(dest)

	if os.IsNotExist(err) {
		return err
	}

	if !info.IsDir() {
		return errors.New("Not a directory")
	}

	method := "cooperhewitt.shoebox.items.getList"

	cb := func(rsp api.APIResponse) error {

		items := gjson.GetBytes(rsp.Raw(), "items")

		for _, i := range items.Array() {

			b := []byte(i.Raw)
			err := sb.ArchiveItem(dest, b)

			if err != nil {
				return err
			}
		}

		return nil
	}

	args := url.Values{}

	err = sb.client.ExecuteMethodPaginated(method, &args, cb)

	if err != nil {
		return err
	}

	return nil
}

func (sb *Shoebox) ArchiveItem(root string, item []byte) error {

	var err error

	item_id := gjson.GetBytes(item, "id").Int()

	rel_path := util.Id2Path(item_id)
	root_path := filepath.Join(root, rel_path)

	_, err = os.Stat(root_path)

	if os.IsNotExist(err) {

		err = os.MkdirAll(root_path, 0755)

		if err != nil {
			return err
		}
	}

	// TO DO
	// generate HTML

	err = sb.ArchiveItemMetadata(root, item)

	if err != nil {
		return err
	}

	err = sb.ArchiveItemObject(root, item)

	if err != nil {
		return err
	}

	return nil
}

func (sb *Shoebox) ArchiveItemMetadata(root string, item []byte) error {

	id := gjson.GetBytes(item, "id").Int()
	path := util.Id2Path(id)

	fname := fmt.Sprintf("%d.json", id)

	rel_path := filepath.Join(path, fname)
	abs_path := filepath.Join(root, rel_path)

	err := ioutil.WriteFile(abs_path, item, 0644)

	if err != nil {
		return err
	}

	log.Println("write", abs_path)
	return nil
}

func (sb *Shoebox) ArchiveItemObject(root string, item []byte) error {

	object_id := gjson.GetBytes(item, "refers_to_uid").Int()

	method := "cooperhewitt.objects.getInfo"
	args := url.Values{}

	str_id := strconv.FormatInt(object_id, 10)
	args.Set("object_id", str_id)

	rsp, err := sb.client.ExecuteMethod(method, &args)

	if err != nil {
		return err
	}

	object := rsp.Raw()

	err = sb.ArchiveItemObjectMetadata(root, item, object)

	if err != nil {
		return err
	}

	err = sb.ArchiveItemObjectImages(root, item, object)

	if err != nil {
		return err
	}

	return nil
}

func (sb *Shoebox) ArchiveItemObjectMetadata(root string, item []byte, object []byte) error {

	item_id := gjson.GetBytes(item, "id").Int()
	object_id := gjson.GetBytes(item, "refers_to_uid").Int()

	path := util.Id2Path(item_id)
	fname := fmt.Sprintf("%d.json", object_id)

	rel_path := filepath.Join(path, fname)
	abs_path := filepath.Join(root, rel_path)

	err := ioutil.WriteFile(abs_path, object, 0644)

	if err != nil {
		return err
	}

	log.Println("write", abs_path)
	return nil
}

func (sb *Shoebox) ArchiveItemObjectImages(root string, item []byte, object []byte) error {

	item_id := gjson.GetBytes(item, "id").Int()
	path := util.Id2Path(item_id)

	rel_path := filepath.Join(root, path)

	images := gjson.GetBytes(object, "object.images")

	t1 := time.Now()

	wg := new(sync.WaitGroup)

	for _, img := range images.Array() {

		for _, details := range img.Map() {

			url := details.Get("url")

			remote := url.String()
			fname := filepath.Base(remote)

			local := filepath.Join(rel_path, fname)

			_, err := os.Stat(local)

			if os.IsExist(err) {
				log.Printf("skip %s\n", remote)
				continue
			}

			wg.Add(1)

			go func(remote string, local string, wg *sync.WaitGroup) {

				defer wg.Done()

				util.GetStore(remote, local)
				log.Println("write", local)				

			}(remote, local, wg)

		}
	}

	wg.Wait()

	t2 := time.Since(t1)

	log.Printf("time to get images for %d : %v\n", item_id, t2)
	return nil
}
