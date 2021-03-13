# golang-ifood-sdk
A golang Ifood sdk 


## Usage



```go

package main

import (
    sdk "github.com/arxdsilva/golang-ifood-sdk/container"
)

func main() {
    var clientID, clientSecret, user, password string
    clientID = os.GetEnv("CLIENT_ID")
    clientSecret = os.GetEnv("CLIENT_SECRET")
    container := sdk.New(0, time.Minute)
    container.GetHttpAdapter()
    auth := container.GetAuthenticationService(clientID, clientSecret)
    user = os.GetEnv("USER")
    password = os.GetEnv("PASSWORD")
    creds, err := auth.Authenticate(user,password)
    if err != nil { 
        log.Fatal(err) 
    }
    merchant := container.GetMerchantService(creds.AccessToken)
    merchants, err := merchant.ListAll()
    if err != nil { 
        log.Fatal(err) 
    }
    fmt.Printf("creds: %+v\n", creds)
    fmt.Printf("merchants: %+v\n", merchants)
}
```