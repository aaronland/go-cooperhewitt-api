package shoebox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-cooperhewitt-api"
	"github.com/thisisaaronland/go-cooperhewitt-api/util"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type Shoebox struct {
	client api.APIClient
}

type CHShoeboxItems struct {
	Items []*CHShoeboxItem `json:"items"`
}

type CHShoeboxItem struct {
	Id           string            `json:"id"`
	Created      string            `json:"created"`
	LastModified string            `json:"lastmodified"`
	IsPublic     string            `json:"is_public"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	RefersToId   string            `json:"refers_to_id"`
	RefersToUid  string            `json:"refers_to_uid"`
	RefersToA    string            `json:"refers_to_a"`
	RefersTo     CHShoeboxRefersTo `json:"refers_to"`
	Action       string            `json:"action"`
}

type CHShoeboxRefersTo struct {
	AccessionNumber     string `json:"accession_number"`
	CreditLine          string `json:"credit_line"`
	Date                string `json:"date"`
	Decade              string `json:"decade"`
	DepartmentId        string `json:"department_id"`
	Description         string `json:"description"`
	Dimensions          string `json:"dimentions"`
	DimensionsRaw       string `json:"dimentions_raw"`
	GalleryText         string `json:"gallery_text"`
	HasNoKnownCopyright string `json:"has_no_known_copyright"`
	Id                  string `json:"id"`
	Inscribed           string `json:"inscribed"`
	IsLoanObject        int    `json:"is_loan_object"`
	Justification       string `json:"justification"`
	LabelText           string `json:"label_text"`
	Markings            string `json:"markings"`
	MediaId             string `json:"media_id"`
	Medium              string `json:"medium"`
	OnDisplay           string `json:"on_display"`
	PeriodId            string `json:"period_id"`
	Provenance          string `json:"provenance"`
	Signed              string `json:"signed"`
	Title               string `json:"title"`
	TitleRaw            string `json:"title_raw"`
	TMSId               string `json:"tms_id"`
	TypeId              string `json:"type_id"`
	URL                 string `json:"url"`
	Videos              string `json:"videos"`
	WOECountry          string `json:"woe:country"`
	WOECountryId        string `json:"woe:country_id"`
	YearAcquired        string `json:"year_acquired"`
	YearEnd             string `json:"year_end"`
	YearStart           string `json:"year_start"`
}

type ShoeboxItems struct {
	Items []*ShoeboxItem `json:"items"`
}

type ShoeboxItem struct {
	Id           int64            `json:"id"`
	Created      int              `json:"created"`
	LastModified int              `json:"lastmodified"`
	IsPublic     bool             `json:"is_public"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	RefersToId   int              `json:"refers_to_id"`
	RefersToUid  int64            `json:"refers_to_uid"`
	RefersToA    string           `json:"refers_to_a"`
	RefersTo     *ShoeboxRefersTo `json:"refers_to"`
	Action       string           `json:"action"`
}

type ShoeboxRefersTo struct {
	AccessionNumber     string `json:"accession_number"`
	CreditLine          string `json:"credit_line"`
	Date                string `json:"date"`
	Decade              int    `json:"decade"`
	DepartmentId        int64  `json:"department_id"`
	Description         string `json:"description"`
	Dimensions          string `json:"dimentions"`
	DimensionsRaw       string `json:"dimentions_raw"`
	GalleryText         string `json:"gallery_text"`
	HasNoKnownCopyright bool   `json:"has_no_known_copyright"`
	Id                  int64  `json:"id"`
	Inscribed           string `json:"inscribed"`
	IsLoanObject        bool   `json:"is_loan_object"`
	Justification       string `json:"justification"`
	LabelText           string `json:"label_text"`
	Markings            string `json:"markings"`
	MediaId             int64  `json:"media_id"`
	Medium              string `json:"medium"`
	OnDisplay           bool   `json:"on_display"`
	PeriodId            int64  `json:"period_id"`
	Provenance          string `json:"provenance"`
	Signed              string `json:"signed"`
	Title               string `json:"title"`
	TitleRaw            string `json:"title_raw"`
	TMSId               int    `json:"tms_id"`
	TypeId              int64  `json:"type_id"`
	URL                 string `json:"url"`
	Videos              string `json:"videos"`
	WOECountry          int    `json:"woe:country"`
	WOECountryID        int    `json:"woe:country_id"`
	YearAcquired        int    `json:"year_acquired"`
	YearEnd             int    `json:"year_end"`
	YearStart           int    `json:"year_start"`
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

	items, err := sb.GetItems()

	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)

	for _, i := range items.Items {

		wg.Add(1)

		go func(i *ShoeboxItem, wg *sync.WaitGroup) {

			defer wg.Done()

			err := sb.ArchiveItem(dest, i)

			if err != nil {
				log.Println(err)
			}

		}(i, wg)
	}

	wg.Wait()

	return nil
}

func (sb *Shoebox) ArchiveItem(root string, item *ShoeboxItem) error {

	// TO DO
	// fetch object info
	// fetch images (videos, etc)
	// generate HTML

	body, err := json.Marshal(item)

	if err != nil {
		return err
	}

	id := item.Id
	path := util.Id2Path(id)

	fname := fmt.Sprintf("%d.json", id)

	rel_path := filepath.Join(path, fname)
	abs_path := filepath.Join(root, rel_path)

	root_path := filepath.Dir(abs_path)

	_, err = os.Stat(root_path)

	if os.IsNotExist(err) {

		err = os.MkdirAll(root_path, 0755)

		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(abs_path, body, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (sb *Shoebox) GetItems() (*ShoeboxItems, error) {

	items := make([]*ShoeboxItem, 0)

	method := "cooperhewitt.shoebox.items.getList"

	cb := func(rsp api.APIResponse) error {

		var ch_items CHShoeboxItems
		err := json.Unmarshal(rsp.Raw(), &ch_items)

		if err != nil {
			return err
		}

		for _, ch_item := range ch_items.Items {

			sb_item, err := CHItemToSBItem(ch_item)

			if err != nil {
				return err
			}

			items = append(items, sb_item)
		}

		return nil
	}

	args := url.Values{}

	err := sb.client.ExecuteMethodPaginated(method, &args, cb)

	if err != nil {
		return nil, err
	}

	return &ShoeboxItems{items}, nil
}

func CHItemToSBItem(ch_item *CHShoeboxItem) (*ShoeboxItem, error) {

	id, err := strconv.ParseInt(ch_item.Id, 10, 64)

	if err != nil {
		return nil, err
	}

	created, err := strconv.Atoi(ch_item.Created)

	if err != nil {
		return nil, err
	}

	lastmod, err := strconv.Atoi(ch_item.LastModified)

	if err != nil {
		return nil, err
	}

	refersto_id, err := strconv.Atoi(ch_item.RefersToId)

	if err != nil {
		return nil, err
	}

	refersto_uid, err := strconv.ParseInt(ch_item.RefersToUid, 10, 64)

	if err != nil {
		return nil, err
	}

	is_public := false

	if ch_item.IsPublic == "1" {
		is_public = true
	}

	sb_refersto, err := CHRefersToToSBRefersTo(&ch_item.RefersTo)

	if err != nil {
		return nil, err
	}

	sb_item := ShoeboxItem{
		Id:           id,
		Created:      created,
		LastModified: lastmod,
		Title:        ch_item.Title,
		Description:  ch_item.Description,
		Action:       ch_item.Action,
		RefersTo:     sb_refersto,
		RefersToA:    ch_item.RefersToA,
		RefersToId:   refersto_id,
		RefersToUid:  refersto_uid,
		IsPublic:     is_public,
	}

	return &sb_item, nil
}

func CHRefersToToSBRefersTo(ch_refersto *CHShoeboxRefersTo) (*ShoeboxRefersTo, error) {

	decade, err := strconv.Atoi(ch_refersto.Decade)

	if err != nil {
		decade = 0
	}

	department_id, err := strconv.ParseInt(ch_refersto.DepartmentId, 10, 64)

	if err != nil {
		department_id = 0
	}

	copyright := true

	if ch_refersto.HasNoKnownCopyright == "1" {
		copyright = false
	}

	id, err := strconv.ParseInt(ch_refersto.Id, 10, 64)

	if err != nil {
		id = 0
	}

	loan_object := true

	if ch_refersto.IsLoanObject == 1 {
		loan_object = false
	}

	media_id, err := strconv.ParseInt(ch_refersto.MediaId, 10, 64)

	if err != nil {
		media_id = 0
	}

	on_display := false

	if ch_refersto.OnDisplay == "1" {
		loan_object = true
	}

	period_id, err := strconv.ParseInt(ch_refersto.PeriodId, 10, 64)

	if err != nil {
		period_id = 0
	}

	tms_id, err := strconv.Atoi(ch_refersto.TMSId)

	if err != nil {
		tms_id = 0
	}

	type_id, err := strconv.ParseInt(ch_refersto.TypeId, 10, 64)

	if err != nil {
		type_id = 0
	}

	woe_country, err := strconv.Atoi(ch_refersto.WOECountry)

	if err != nil {
		woe_country = 0
	}

	woe_country_id, err := strconv.Atoi(ch_refersto.WOECountryId)

	if err != nil {
		woe_country_id = 0
	}

	year_acquired, err := strconv.Atoi(ch_refersto.YearAcquired)

	if err != nil {
		year_acquired = 0
	}

	year_end, err := strconv.Atoi(ch_refersto.YearEnd)

	if err != nil {
		year_end = 0
	}

	year_start, err := strconv.Atoi(ch_refersto.YearStart)

	if err != nil {
		year_start = 0
	}

	rf := ShoeboxRefersTo{
		AccessionNumber:     ch_refersto.AccessionNumber,
		CreditLine:          ch_refersto.CreditLine,
		Date:                ch_refersto.Date,
		Decade:              decade,
		DepartmentId:        department_id,
		Description:         ch_refersto.Description,
		Dimensions:          ch_refersto.Dimensions,
		DimensionsRaw:       ch_refersto.DimensionsRaw,
		GalleryText:         ch_refersto.GalleryText,
		HasNoKnownCopyright: copyright,
		Id:                  id,
		Inscribed:           ch_refersto.Inscribed,
		IsLoanObject:        loan_object,
		Justification:       ch_refersto.Justification,
		LabelText:           ch_refersto.LabelText,
		Markings:            ch_refersto.Markings,
		MediaId:             media_id,
		Medium:              ch_refersto.Medium,
		OnDisplay:           on_display,
		PeriodId:            period_id,
		Provenance:          ch_refersto.Provenance,
		Signed:              ch_refersto.Signed,
		Title:               ch_refersto.Title,
		TitleRaw:            ch_refersto.TitleRaw,
		TMSId:               tms_id,
		TypeId:              type_id,
		URL:                 ch_refersto.URL,
		Videos:              ch_refersto.Videos,
		WOECountry:          woe_country,
		WOECountryID:        woe_country_id,
		YearAcquired:        year_acquired,
		YearEnd:             year_end,
		YearStart:           year_start,
	}

	return &rf, nil
}
