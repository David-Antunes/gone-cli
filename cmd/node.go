package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/David-Antunes/gone/api/Add"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manages nodes in the network emulation",
	Long:  `Node command starts a container in the emulation`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing docker command")
			return
		}
		machine, _ := cmd.Flags().GetString("machine")
		dockerCmd := strings.Join(args, " ")
		fmt.Println(dockerCmd)
		if args[0] != "docker" {
			fmt.Println("Only docker is supported")
		}
		client := startClient()

		body, err := json.Marshal(&api.AddNodeRequest{
			DockerCmd: dockerCmd,
			MachineId: machine,
		})

		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", url+"/addNode", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		d := json.NewDecoder(res.Body)
		resp := api.AddNodeResponse{}
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
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringP("machine", "m", "", "Starts docker container in the provided machine id")
}
