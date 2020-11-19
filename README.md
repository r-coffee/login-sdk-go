# log in service SDK

### Example Usage
```golang
package main

import (
	"fmt"
	"log"

	login "github.com/r-coffee/login-sdk-go"
)

func main() {
	// create a new login client with a valid entity id
	client := login.CreateClient("Eg4KBmVudGl0eRoEVEVTVA")

	// test register
	err := client.Register("foo", "password")
	if err != nil {
		log.Fatal(err)
	}

	// test login
	token, err := client.Login("foo", "password")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	// test validate
	err = client.Validate(token)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}
```