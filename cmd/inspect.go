package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/David-Antunes/gone/api/Inspect"
	"net/http"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Args:  cobra.ExactArgs(1),
	Short: "Shows information about the emulation",
	Long:  `Shows more information about the emulation`,
	Run: func(cmd *cobra.Command, args []string) {

		object := args[0]

		node, _ := cmd.Flags().GetBool("node")

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		client := startClient()

		var req *http.Request

		if node {
			body, err := json.Marshal(&api.InspectNodeRequest{Name: object})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/inspectNode", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			req.Header.Add("Content-Type", "application/json")

			jsonOutput(cmd, body, req)
			res, err := client.Do(req)
			d := json.NewDecoder(res.Body)
			resp := api.InspectNodeResponse{}

			err = d.Decode(&resp)

			if err != nil {
				panic(err)
			}
			if resp.Error.ErrCode != 0 {
				fmt.Println(resp.Error.ErrMsg)
				return
			}
			jsonOut, err := json.MarshalIndent(resp, "", "\t")

			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonOut))

		} else if bridge {

			body, err := json.Marshal(&api.InspectBridgeRequest{Name: object})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/inspectBridge", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			req.Header.Add("Content-Type", "application/json")

			jsonOutput(cmd, body, req)
			res, err := client.Do(req)
			d := json.NewDecoder(res.Body)
			resp := api.InspectBridgeResponse{}

			err = d.Decode(&resp)

			if err != nil {
				panic(err)
			}

			if resp.Error.ErrCode != 0 {
				fmt.Println(resp.Error.ErrMsg)
				return
			}

			jsonOut, err := json.MarshalIndent(resp, "", "\t")

			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonOut))

		} else if router {

			body, err := json.Marshal(&api.InspectRouterRequest{Name: object})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/inspectRouter", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			req.Header.Add("Content-Type", "application/json")

			jsonOutput(cmd, body, req)
			res, err := client.Do(req)

			d := json.NewDecoder(res.Body)
			resp := api.InspectRouterResponse{}

			err = d.Decode(&resp)

			if err != nil {
				panic(err)
			}
			if resp.Error.ErrCode != 0 {
				fmt.Println(resp.Error.ErrMsg)
				return
			}

			prettyPrint(resp)

		} else {
			fmt.Println()
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolP("node", "n", false, "Shows node information")
	inspectCmd.Flags().BoolP("bridge", "b", false, "Shows bridge information")
	inspectCmd.Flags().BoolP("router", "r", false, "Shows router information")
	inspectCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router")

}
