package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"moul.io/http2curl/v2"
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
	if b, _ := rootCmd.Flags().GetBool("quiet"); !b {
		fmt.Println(string(jsonOut))
	}
}

func jsonOutput(cmd *cobra.Command, body []byte, req *http.Request) {
	//fmt.Println(cmd.Flags().GetBool("json"))
	if b, _ := cmd.Flags().GetBool("json"); b {
		fmt.Println(string(body))
	}

	if c, _ := cmd.Flags().GetBool("curl"); c {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Println(command)
	}

	if b, _ := cmd.Flags().GetBool("dry-run"); b {
		os.Exit(0)
	}
}
