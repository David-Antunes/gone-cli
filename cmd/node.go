package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/David-Antunes/gone/api/Add"
	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use: "node [flags] [-m machineId] -- <docker command>",
	Example: `
	gone-cli node -m machine1 -- docker run -d --network gone_net --name ubuntu1 ubuntu sleep 10000 

Adds the container to the emulation in the instance machine1. The docker command
requires the -d and --network gone_net to sucessfully execute.`,
	Short: "Adds node to the network emulation",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing docker command")
			return
		}
		// fmt.Println(args)
		// dockerCmd := ""
		// envVar := false
		// for _, arg := range args {
		// 	if envVar {
		// 		env := strings.Split(arg, "=")
		// 		if len(env) == 2 {
		// 			nextArg := strings.Split(env[1], " ")
		// 			if len(nextArg) >= 2 {
		// 				newS := strings.Join(env[1:], " ")
		// 				dockerCmd = dockerCmd + env[0] + "=\"" + newS + "\""
		// 				envVar = false
		// 				continue
		// 			} else {
		// 				dockerCmd = dockerCmd + " " + env[0] + "=" + env[1]
		// 				envVar = false
		// 				continue
		// 			}
		// 		} else {
		// 				newS := strings.Join(env[1:], "=")
		// 				dockerCmd = dockerCmd + " " + env[0] + "=\"" + newS + "\""
		// 				envVar = false
		// 				continue
		// 		}
		// 	}
		// 	if arg == "-e" {
		// 		envVar = true
		// 	}
		// 	dockerCmd = dockerCmd + " " + arg
		// }

		machine, _ := cmd.Flags().GetString("machine")
		// dockerCmd = strings.Trim(dockerCmd, " ")
		// fmt.Println(dockerCmd)
		if args[0] != "docker" {
			fmt.Println("Only docker is supported")
			return
		}
		client := startClient()

		body, err := json.Marshal(&api.AddNodeRequest{
			DockerCmd: args,
			MachineId: machine,
		})

		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", URL+"/addNode", bytes.NewBuffer(body))
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
