package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var URL = "http://localhost:3000"
var AGENT_URL = "http://localhost:3300"

func startClient() *http.Client {
	return http.DefaultClient
}

func prettyPrint(data any) {
	jsonOut, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonOut))
}
