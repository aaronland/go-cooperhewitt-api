package main

import (
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-cooperhewitt-api"
	"github.com/thisisaaronland/go-cooperhewitt-api/client"
	"github.com/thisisaaronland/go-cooperhewitt-api/endpoint"
	"log"
	"os"
)

func main() {

	var api_params api.APIParams

	flag.Var(&api_params, "param", "One or more Who's On First API query=value parameters.")

	// var stdout = flag.Bool("stdout", false, "Write API results to STDOUT")
	// var raw = flag.Bool("raw", false, "Dump raw Who's On First API responses.")
	// var paginated = flag.Bool("paginated", false, "Automatically paginate API results.")

	var custom_endpoint = flag.String("endpoint", "", "Define a custom endpoint for the Cooper Hewitt API.")

	flag.Parse()

	args := api_params.ToArgs()

	token := args.Get("access_token")
	method := args.Get("method")

	if method == "" {
		log.Fatal("You forgot to specify a method")
	}

	e, err := endpoint.NewOAuth2APIEndpoint(token)

	if err != nil {
		log.Fatal(err)
	}

	if *custom_endpoint != "" {

		err := e.SetEndpoint(*custom_endpoint)

		if err != nil {
			log.Fatal(err)
		}

	}

	c, _ := client.NewHTTPClient(e)

		dest := os.Stdout

		cb := func(rsp api.APIResponse) error {
			_, err = dest.Write(rsp.Raw())
			return err
		}

	err = c.ExecuteMethodWithCallback(method, args, cb)

	if err != nil {
		msg := fmt.Sprintf("Failed to call '%s' because %s", method, err)
		log.Fatal(msg)
	}

	os.Exit(0)
}
