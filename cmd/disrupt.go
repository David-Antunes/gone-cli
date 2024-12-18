package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/David-Antunes/gone/api/Operations"
	"github.com/spf13/cobra"
	"net/http"
)

var disruptCmd = &cobra.Command{
	Use:   "disrupt",
	Args:  cobra.MinimumNArgs(1),
	Short: "Disrupts network components in the network emulation",
	Run: func(cmd *cobra.Command, args []string) {

		client := startClient()

		firstObject := args[0]
		var secondObject string
		node, _ := cmd.Flags().GetBool("node")

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		stop, _ := cmd.Flags().GetBool("stop")

		if router {
			if len(args) == 2 {
				secondObject = args[1]
			} else {
				fmt.Println("Missing second router")
				return
			}
		}

		var req *http.Request
		var disruptURL string
		var body []byte
		var err error
		if node {

			if !stop {
				disruptURL = "/disruptNode"
			} else {
				disruptURL = "/stopDisruptNode"
			}
			body, err = json.Marshal(&api.DisruptNodeRequest{
				Node: firstObject,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+disruptURL, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if bridge {
			if !stop {
				disruptURL = "/disruptBridge"
			} else {
				disruptURL = "/stopDisruptBridge"
			}
			body, err = json.Marshal(&api.DisruptBridgeRequest{
				Bridge: firstObject,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+disruptURL, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if router {
			if !stop {
				disruptURL = "/disruptRouters"
			} else {
				disruptURL = "/stopDisruptRouters"
			}
			body, err = json.Marshal(&api.DisruptRoutersRequest{
				Router1: firstObject,
				Router2: secondObject,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+disruptURL, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Missing -n/-b/-r flag")
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
			resp = api.DisruptNodeResponse{}
		} else if bridge {
			resp = api.DisruptBridgeResponse{}
		} else if router {
			resp = api.DisruptRoutersResponse{}
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
	rootCmd.AddCommand(disruptCmd)
	disruptCmd.Flags().BoolP("node", "n", false, "Disrupts link between node and bridge")
	disruptCmd.Flags().BoolP("bridge", "b", false, "Disrupts link between bridge and router")
	disruptCmd.Flags().BoolP("router", "r", false, "Disrupts link between routers")
	disruptCmd.Flags().BoolP("stop", "s", false, "Stops Disruption")
	disruptCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router")
}
