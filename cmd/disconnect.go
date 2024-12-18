/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/David-Antunes/gone/api/Disconnect"
	"github.com/spf13/cobra"
	"net/http"
)

// disconnectCmd represents the disconnect command
var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnects components from the emulation",
	Long:  ``,
	Args:  cobra.MaximumNArgs(2),
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
			fmt.Println(err)
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

	disconnectCmd.Flags().BoolP("node", "n", false, "Disconnects node")
	disconnectCmd.Flags().BoolP("bridge", "b", false, "Disconnects bridge")
	disconnectCmd.Flags().BoolP("router", "r", false, "Disconnects router")
	disconnectCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router")
}
