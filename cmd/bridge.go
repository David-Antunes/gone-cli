package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/David-Antunes/gone/api/Add"

	"github.com/spf13/cobra"
)

// bridgeCmd represents the bridge command
var bridgeCmd = &cobra.Command{
	Use: "bridge [flags] [-m machineId] <id>",
	Example: `
	gone-cli bridge bridge1

Creates bridge1 in Leader instance

	gone-cli bridge -m machine1 bridge1

Creates bridge1 in instance machine1`,
	Args:  cobra.ExactArgs(1),
	Short: "Adds a bridge to emulation",
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		machineId, _ := cmd.Flags().GetString("machine")
		client := startClient()

		body, err := json.Marshal(&api.AddBridgeRequest{
			Name:      id,
			MachineId: machineId,
		})

		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", URL+"/addBridge", bytes.NewBuffer(body))
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
		resp := api.AddBridgeResponse{}
		err = d.Decode(&resp)

		if err != nil {
			panic(err)
		}
		if resp.Error.ErrCode != 0 {
			fmt.Println(resp.Error.ErrMsg)
			return
		}
		prettyPrint(resp)

	},
}

func init() {
	rootCmd.AddCommand(bridgeCmd)
	bridgeCmd.Flags().StringP("machine", "m", "", "Starts a bridge in the provided machine id")
}
