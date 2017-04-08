package util

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
)

func HTTPResponseToReader(http_rsp *http.Response) (io.Reader, error) {

	var body io.Reader
	var err error

	switch http_rsp.Header.Get("Content-Encoding") {

	case "gzip":

		body, err = gzip.NewReader(http_rsp.Body)

		if err != nil {
			return nil, err
		}

	default:
		body = http_rsp.Body
	}

	return body, nil
}

func HTTPResponseToBytes(http_rsp *http.Response) ([]byte, error) {

	body, err := HTTPResponseToReader(http_rsp)

	if err != nil {
		return nil, err
	}

	http_body, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	return http_body, nil
}
