# go-cooperhewitt-api

Too soon. Move along.

## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

### Fancy

```
import (
	"github.com/thisisaaronland/go-cooperhewitt-api"
	"github.com/thisisaaronland/go-cooperhewitt-api/client"
	"github.com/tidwall/gjson"
	"log"	
)

func main() {

     	tk := "ACCESS_TOKEN"
	
	ep, _ := endpoint.NewOAuth2APIEndpoint(tk)
	cl, _ := client.NewHTTPClient(ep)

	method := "cooperhewitt.shoebox.items.getList"

	cb := func(rsp api.APIResponse) error {

		items := gjson.GetBytes(rsp.Raw(), "items")
		
		for _, i := range items.Array() {

			item := []byte(i.Raw)
			item_id := gjson.GetBytes(item, "id").Int()

			log.Println(item_id)			
		}

		return nil
	}

	args := url.Values{}

	cl.ExecuteMethodPaginated(method, &args, cb)
}	
```

## Tools

### ch-api

```
./bin/ch-api -param access_token=TOKEN -param method=cooperhewitt.labs.whatWouldMicahSay
{
  "micah": {
    "says": "Can you direct me to the shuffle board deck?"
  }, 
  "stat": "ok"
}
```
