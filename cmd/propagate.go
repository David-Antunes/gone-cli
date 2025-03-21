/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	api "github.com/David-Antunes/gone/api/Operations"

	"github.com/spf13/cobra"
)

// propagateCmd represents the propagate command
var propagateCmd = &cobra.Command{
	Use: "propagate [flags] <id>",
	Example: `
	gone-cli propagate router1
	
Propagates router1's routing rules to adjacent routers`,
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := args[0]

		client := startClient()

		body, err := json.Marshal(&api.PropagateRequest{
			Name: id,
		})

		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", URL+"/propagate", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")

		jsonOutput(cmd, body, req)
		res, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		d := json.NewDecoder(res.Body)
		resp := api.PropagateResponse{}
		err = d.Decode(&resp)

		if err != nil {
			panic(err)
		}
		prettyPrint(resp)
	},
}

func init() {
	rootCmd.AddCommand(propagateCmd)
}
