package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	api "github.com/David-Antunes/gone/api/Add"

	"github.com/spf13/cobra"
)

// routerCmd represents the router command

var routerCmd = &cobra.Command{
	Use: "router [flags] [-m machineId] <id>",
	Example: `
	gone-cli router -m machine1 router1

Creates router1 in instance machine1`,
	Args:  cobra.ExactArgs(1),
	Short: "Adds new router to the network emulation.",
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		client := startClient()

		machineId, _ := cmd.Flags().GetString("machine")
		body, err := json.Marshal(&api.AddRouterRequest{
			Name:      id,
			MachineId: machineId,
		})
		//fmt.Println(string(body))
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", URL+"/addRouter", bytes.NewBuffer(body))
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
		resp := api.AddRouterResponse{}
		err = d.Decode(&resp)

		if err != nil {
			panic(err)
		}
		prettyPrint(resp)

	},
}

func init() {
	rootCmd.AddCommand(routerCmd)
	routerCmd.Flags().StringP("machine", "m", "", "Starts a router in the provided machine id")
}
