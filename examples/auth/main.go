package main

import (
	"fmt"
	"log"
	"os"

	sdk "github.com/arxdsilva/golang-ifood-sdk/container"
)

func main() {
	var clientID, clientSecret string
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	// START new SDK instance
	c := sdk.Create(clientID, clientSecret, 0, true)
	// Get user code to connect this supplier to the restaurant
	uc, err := c.AuthService.V2GetUserCode()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("v2 user: %+v\n", uc)
	v2Creds, err := c.AuthService.V2Authenticate(
		"client_credentials",
		uc.Usercode, uc.AuthorizationCodeVerifier, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("v2 creds: %+v\n", v2Creds)
}
