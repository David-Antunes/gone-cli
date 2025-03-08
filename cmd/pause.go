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

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use: "pause [flags] [-a] <id>",
	Example: `
	gone-cli pause node1

Pauses node1's execution

	gone-cli pause -a

Pauses all nodes`,
	Short: "Pauses a node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := ""
		if len(args) == 1 {
			id = args[0]
		}

		all, _ := cmd.Flags().GetBool("all")

		if id == "" && !all {
			cmd.Help()
			return
		}
		client := startClient()

		body, err := json.Marshal(&api.PauseRequest{
			Id:  id,
			All: all,
		})
		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", URL+"/pause", bytes.NewBuffer(body))
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
		resp := api.PauseResponse{}
		err = d.Decode(&resp)

		if err != nil {
			panic(err)
		}
		prettyPrint(resp)
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
	pauseCmd.Flags().BoolP("all", "a", false, "Pauses all nodes")
}
