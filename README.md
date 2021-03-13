# golang-ifood-sdk
A golang Ifood sdk 


## Usage



```go

package main

import (
    sdk "github.com/arxdsilva/golang-ifood-sdk/container"
)

func main() {
    var clientID, clientSecret string
    clientID = os.GetEnv("CLIENT_ID")
    clientSecret = os.GetEnv("CLIENT_SECRET")
    container := sdk.New(0, time.Minute)
    container.GetHttpAdapter()
    auth := container.GetAuthenticationService(clientID, clientSecret)
    creds, err := auth.Authenticate(user,password)
    // err check
    var token string
    merchant := container.GetMerchantService(creds.AccessToken)
    merchants, err := merchant.ListAll()
    // err check
}
```