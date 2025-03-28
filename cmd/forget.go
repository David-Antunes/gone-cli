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

// forgetCmd represents the forget command
var forgetCmd = &cobra.Command{
	Use: "forget [flags] <id>",
	Example: `
	gone-cli forget router1

Cleans the routing table of router1`,
	Short: "Cleans routing rules for a specific router",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := args[0]

		client := startClient()

		body, err := json.Marshal(&api.ForgetRequest{
			Name: id,
		})

		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", URL+"/forget", bytes.NewBuffer(body))
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
		resp := api.ForgetResponse{}
		err = d.Decode(&resp)

		if err != nil {
			panic(err)
		}
		prettyPrint(resp)

	},
}

func init() {
	rootCmd.AddCommand(forgetCmd)
}
