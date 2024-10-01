/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	api "github.com/David-Antunes/gone/api/Operations"
	"net/http"

	"github.com/spf13/cobra"
)

// unpauseCmd represents the unpause command
var unpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "A brief description of your command",
	Long:  ``,
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

		body, err := json.Marshal(&api.UnpauseRequest{
			Id:  id,
			All: all,
		})

		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", url+"/unpause", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		d := json.NewDecoder(res.Body)
		resp := api.UnpauseResponse{}
		err = d.Decode(&resp)

		if err != nil {
			panic(err)
		}
		prettyPrint(resp)
	},
}

func init() {
	rootCmd.AddCommand(unpauseCmd)
	unpauseCmd.Flags().BoolP("all", "a", false, "Pauses all nodes")
}
