/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/David-Antunes/gone/api/Remove"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use: "remove [flags] {-n | -b | -r} <id>",
	Example: `
	gone-cli remove -n node1

Removes node1 from emulation
	
	gone-cli remove -r router1

Removes router1 from emulation`,
	Args:  cobra.ExactArgs(1),
	Short: "Removes a component from the network emulation",
	Long:  `Removes a component from the network emulation`,
	Run: func(cmd *cobra.Command, args []string) {

		object := args[0]

		node, _ := cmd.Flags().GetBool("node")

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		client := startClient()

		var req *http.Request
		var body []byte
		var err error

		if node {
			body, err = json.Marshal(&api.RemoveNodeRequest{Name: object})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/removeNode", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if bridge {

			body, err = json.Marshal(&api.RemoveBridgeRequest{Name: object})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/removeBridge", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if router {

			body, err = json.Marshal(&api.RemoveRouterRequest{Name: object})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/removeRouter", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Missing appropriate flag")
			return
		}

		req.Header.Add("Content-Type", "application/json")
		jsonOutput(cmd, body, req)
		res, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		d := json.NewDecoder(res.Body)
		var resp any
		if node {
			resp = api.RemoveNodeResponse{}
		} else if bridge {
			resp = api.RemoveBridgeResponse{}
		} else if router {
			resp = api.RemoveRouterResponse{}
		} else {
			return
		}

		err = d.Decode(&resp)

		if err != nil {
			panic(err)
		}
		prettyPrint(resp)

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("node", "n", false, "Removes node")
	removeCmd.Flags().BoolP("bridge", "b", false, "Removes bridge")
	removeCmd.Flags().BoolP("router", "r", false, "Removes router")
	removeCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router")
}
