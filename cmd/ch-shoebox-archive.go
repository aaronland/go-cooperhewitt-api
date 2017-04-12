package main

import (
	"flag"
	"github.com/briandowns/spinner"
	"github.com/thisisaaronland/go-cooperhewitt-api/client"
	"github.com/thisisaaronland/go-cooperhewitt-api/endpoint"
	"github.com/thisisaaronland/go-cooperhewitt-api/shoebox"
	"log"
	"os"
	"time"
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

	sb, err := shoebox.NewShoeboxArchiver(c)

	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {

		sp := spinner.New(spinner.CharSets[38], 200*time.Millisecond)
		sp.Prefix = "archiving shoebox..."
		sp.Start()

		for {

			select {
			case <-done:
				sp.Stop()
				return
			}
		}
	}()

	err = sb.Archive(*dest)

	if err != nil {
		log.Fatal(err)
	}

	done <- true

	os.Exit(0)
}
