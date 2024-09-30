package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var unixSock = "/tmp/emu.sock"

var service = "example"

// var unixUrl = "http+unix://" + service
var url = "http://localhost:3000"

//func startClient() *http.Client {
//	u := &httpunix.Transport{
//		DialTimeout:           1000 * time.Millisecond,
//		RequestTimeout:        1 * time.Second,
//		ResponseHeaderTimeout: 1 * time.Second,
//	}
//	u.RegisterLocation(service, unixSock)
//
//	t := &http.Transport{}
//	t.RegisterProtocol(httpunix.Scheme, u)
//
//	return &http.Client{Transport: t}
//}

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
