package shoebox

import (
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-cooperhewitt-api"
	"github.com/thisisaaronland/go-cooperhewitt-api/util"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	_ "time"
)

type ShoeboxArchiver struct {
	client api.APIClient
}

func NewShoeboxArchiver(client api.APIClient) (*ShoeboxArchiver, error) {

	sb := ShoeboxArchiver{
		client: client,
	}

	return &sb, nil
}

func (sb *ShoeboxArchiver) Archive(root_path string) error {

	info, err := os.Stat(root_path)

	if os.IsNotExist(err) {
		return err
	}

	if !info.IsDir() {
		return errors.New("Not a directory")
	}

	max := 2
	throttle := make(chan bool, max)

	for i := 0; i < max; i++ {
		throttle <- true
	}

	method := "cooperhewitt.shoebox.items.getList"

	cb := func(rsp api.APIResponse) error {

		items := gjson.GetBytes(rsp.Raw(), "items")
		wg := new(sync.WaitGroup)

		for _, i := range items.Array() {

			<-throttle

			item := []byte(i.Raw)
			wg.Add(1)

			go func(item []byte, wg *sync.WaitGroup, throttle chan bool) {

				err := sb.ArchiveItem(root_path, item)

				if err != nil {
					log.Println(err)
				}

				throttle <- true

				wg.Done()

			}(item, wg, throttle)

		}

		wg.Wait()
		return nil
	}

	args := url.Values{}

	err = sb.client.ExecuteMethodPaginated(method, &args, cb)

	if err != nil {
		return err
	}

	return nil
}

func (sb *ShoeboxArchiver) ArchiveItem(root string, item []byte) error {

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

	err = sb.ArchiveItemMetadata(root, item)

	if err != nil {
		return err
	}

	// TODO : check for video and other things...
	// (20170410/thisisaaronland)

	var object []byte

	refers_to := gjson.GetBytes(item, "refers_to_a").String()

	switch refers_to {

	case "object":

		object, err = sb.GetItemObject(item)

		if err != nil {
			return err
		}

		err = sb.ArchiveItemObject(root, item, object)

		if err != nil {
			return err
		}

	default:
		log.Printf("TO DO: implement archiving of %s thingies\n", refers_to)
	}

	// TO DO
	// generate HTML

	return nil
}

func (sb *ShoeboxArchiver) ArchiveItemMetadata(root string, item []byte) error {

	id := gjson.GetBytes(item, "id").Int()
	path := util.Id2Path(id)

	fname := "index.json"

	rel_path := filepath.Join(path, fname)
	abs_path := filepath.Join(root, rel_path)

	_, err := os.Stat(abs_path)

	if err == nil {
		return nil
	}

	item_fmt := pretty.Pretty(item)
	err = ioutil.WriteFile(abs_path, item_fmt, 0644)

	if err != nil {
		return err
	}

	log.Println("write", abs_path)
	return nil
}

func (sb *ShoeboxArchiver) GetItemObject(item []byte) ([]byte, error) {

	object_id := gjson.GetBytes(item, "refers_to_uid").Int()

	method := "cooperhewitt.objects.getInfo"
	args := url.Values{}

	str_id := strconv.FormatInt(object_id, 10)
	args.Set("object_id", str_id)

	rsp, err := sb.client.ExecuteMethod(method, &args)

	if err != nil {
		return nil, err
	}

	return rsp.Raw(), nil
}

func (sb *ShoeboxArchiver) ArchiveItemObject(root string, item []byte, object []byte) error {

	var err error

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

func (sb *ShoeboxArchiver) ArchiveItemObjectMetadata(root string, item []byte, object []byte) error {

	item_id := gjson.GetBytes(item, "id").Int()
	object_id := gjson.GetBytes(item, "refers_to_uid").Int()

	path := util.Id2Path(item_id)
	fname := fmt.Sprintf("%d.json", object_id)

	rel_path := filepath.Join(path, fname)
	abs_path := filepath.Join(root, rel_path)

	_, err := os.Stat(abs_path)

	if err == nil {
		return nil
	}

	object_fmt := pretty.Pretty(object)
	err = ioutil.WriteFile(abs_path, object_fmt, 0644)

	if err != nil {
		return err
	}

	log.Println("write", abs_path)
	return nil
}

func (sb *ShoeboxArchiver) ArchiveItemObjectImages(root string, item []byte, object []byte) error {

	item_id := gjson.GetBytes(item, "id").Int()
	path := util.Id2Path(item_id)

	rel_path := filepath.Join(root, path)

	images := gjson.GetBytes(object, "object.images")

	// t1 := time.Now()

	wg := new(sync.WaitGroup)

	for _, img := range images.Array() {

		for _, details := range img.Map() {

			url := details.Get("url")

			remote := url.String()
			fname := filepath.Base(remote)

			local := filepath.Join(rel_path, fname)

			_, err := os.Stat(local)

			if err == nil {
				// log.Printf("skip %s\n", remote)
				continue
			}

			wg.Add(1)

			go func(remote string, local string, wg *sync.WaitGroup) {

				defer wg.Done()

				err := util.GetStore(remote, local)

				if err != nil {
					log.Printf("failed to retrieve %s, because %s\n", remote, err)
					return
				}

				log.Println("write", local)

			}(remote, local, wg)

		}
	}

	wg.Wait()

	// t2 := time.Since(t1)
	// log.Printf("time to get images for %d : %v\n", item_id, t2)

	return nil
}
