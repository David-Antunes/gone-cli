package cmd

import (
	"bytes"
	"encoding/json"
	api "github.com/David-Antunes/gone/api/Add"
	"net/http"

	"github.com/spf13/cobra"
)

// routerCmd represents the router command

var routerCmd = &cobra.Command{
	Use:   "router",
	Args:  cobra.ExactArgs(1),
	Short: "Adds new router to the network emulation.",
	Long:  `Adds new router to the network emulation.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		client := startClient()

		machineId, _ := cmd.Flags().GetString("machine")
		body, err := json.Marshal(&api.AddRouterRequest{
			Name:      id,
			MachineId: machineId,
		})

		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", url+"/addRouter", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		req.Header.Add("Content-Type", "application/json")

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
