# log in service SDK

### Example Usage
```golang
package main

import (
	"fmt"
	"log"

	login "github.com/r-coffee/login-sdk-go"
)

func testSDK() {
	client := login.CreateLoginClient("hostname.com", "my-entity-guid", "/path/to/cert.pem", 8888)

	// test register
	token, err := client.Register("email@example.com", "password")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(token)

	// test login
	token, err = client.Login("email@example.com", "password")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	// test validate
	email, err := client.Validate(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(email)
}

func main() {
	testSDK()
}

```