package main

import (
	"flag"
	"github.com/thisisaaronland/go-cooperhewitt-api/client"
	"github.com/thisisaaronland/go-cooperhewitt-api/endpoint"
	"github.com/thisisaaronland/go-cooperhewitt-api/shoebox"
	"log"
	"os"
)

func main() {

	var token = flag.String("token", "", "")
	var dest = flag.String("dest", "", "")

	flag.Parse()

	e, err := endpoint.NewOAuth2APIEndpoint(*token)

	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewHTTPClient(e)

	if err != nil {
		log.Fatal(err)
	}

	sb, err := shoebox.NewShoebox(c)

	if err != nil {
		log.Fatal(err)
	}

	err = sb.Archive(*dest)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
