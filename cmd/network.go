package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/David-Antunes/gone/api/Operations"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use: "network [flags] [-s] {-b | -r} <id>",
	Example: `
	gone-cli network -s -r router1

Stops router1 from routing packets

	gone-cli network -b bridge1

Starts bridge1`,
	Args:  cobra.ExactArgs(1),
	Short: "Controls whether a particular bridge or router executes",
	Run: func(cmd *cobra.Command, args []string) {

		client := startClient()

		firstObject := args[0]

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		stop, _ := cmd.Flags().GetBool("stop")

		var req *http.Request
		var body []byte
		var err error

		if bridge {
			if !stop {
				body, err = json.MarshalIndent(&api.StartBridgeRequest{
					Bridge: firstObject,
				}, "", "\t")
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/startBridge", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			} else {
				body, err = json.Marshal(&api.StopBridgeRequest{
					Bridge: firstObject,
				})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/stopBridge", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			}
		} else if router {
			if !stop {
				body, err = json.Marshal(&api.StartRouterRequest{
					Router: firstObject,
				})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/startRouter", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			} else {
				body, err = json.Marshal(&api.StopRouterRequest{
					Router: firstObject,
				})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/stopRouter", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println("Missing -b/-r flag")
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

		if stop {
			if bridge {
				resp = api.StopBridgeResponse{}
			} else if router {
				resp = api.StopRouterResponse{}
			} else {
				fmt.Println("Something went wrong")
				return
			}
		} else if bridge {
			resp = api.StartBridgeRequest{}
		} else if router {
			resp = api.StartRouterRequest{}
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
	rootCmd.AddCommand(networkCmd)
	networkCmd.Flags().BoolP("bridge", "b", false, "Bridge to be started or stopped")
	networkCmd.Flags().BoolP("router", "r", false, "Router to be started or stopped")
	networkCmd.Flags().BoolP("stop", "s", false, "Stops component")
	networkCmd.MarkFlagsMutuallyExclusive("bridge", "router")
}
