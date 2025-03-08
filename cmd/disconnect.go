/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/David-Antunes/gone/api/Disconnect"
	"github.com/spf13/cobra"
)

// disconnectCmd represents the disconnect command
var disconnectCmd = &cobra.Command{
	Short: "Disconnects components from the emulation",
	Use:   "disconnect [flags] {-n | -b | -r} <id> [<id>]",
	Example: `
	gone-cli disconnect -n node1

Disconnects node1 from the bridge it is connected to

	gone-cli disconnect -r router1 router2

Disconnects router1 and router2`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			cmd.Help()
			return
		}
		firstObject := args[0]
		var secondObject string
		node, _ := cmd.Flags().GetBool("node")

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		if router {
			if len(args) == 2 {
				secondObject = args[1]
			} else {
				fmt.Println("Missing second router")
				return
			}
		}

		client := startClient()

		var req *http.Request
		var body []byte
		var err error

		if node {
			body, err = json.Marshal(&api.DisconnectNodeRequest{Name: firstObject})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/disconnectNode", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if bridge {

			body, err = json.Marshal(&api.DisconnectBridgeRequest{Name: firstObject})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/disconnectBridge", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if router {

			body, err = json.Marshal(&api.DisconnectRoutersRequest{First: firstObject, Second: secondObject})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/disconnectRouters", bytes.NewBuffer(body))
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
			panic(err)
		}

		d := json.NewDecoder(res.Body)
		var resp any
		if node {
			resp = api.DisconnectNodeResponse{}
		} else if bridge {
			resp = api.DisconnectBridgeResponse{}
		} else if router {
			resp = api.DisconnectRoutersResponse{}
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
	rootCmd.AddCommand(disconnectCmd)

	disconnectCmd.Flags().BoolP("node", "n", false, "Disconnects node from its bridge")
	disconnectCmd.Flags().BoolP("bridge", "b", false, "Disconnects bridge from its router")
	disconnectCmd.Flags().BoolP("router", "r", false, "Disconnects routers")
	disconnectCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router")
}
