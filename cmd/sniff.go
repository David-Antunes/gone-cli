/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/David-Antunes/gone/api/Operations"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// disconnectCmd represents the disconnect command
var sniffCmd = &cobra.Command{
	Use: "sniff [flags] {-n | -b | -r | -s} [-i socket-name] <id> {<id>}",
	Example: `
	gone-cli sniff -i link1 -n node1 

Sniffs network traffic between node1 and its bridge
and redirects it to /tmp/link1.sniff

	gone-cli sniff -s -i link1

Stops sniffing of socket link1

	gone-cli sniff -r router1 router2

Sniffs network traffic between router1 and router2
	
	gone-cli sniff

Shows the list of active sniff links`,
	Short: "Sniffs traffic from Link",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		client := startClient()

		node, _ := cmd.Flags().GetBool("node")

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		stop, _ := cmd.Flags().GetBool("stop")

		id, _ := cmd.Flags().GetString("id")

		if id == "" && !stop {
			u, err := uuid.NewUUID()
			if err != nil {
				panic(err)
			}
			id = u.String()
		} else if id == "" && stop {
			fmt.Println("Missing id")
			return
		}

		if len(args) == 0 && !stop {
			req, err := http.NewRequest("GET", URL+"/listSniffers", nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			d := json.NewDecoder(res.Body)
			list := &api.ListSniffersResponse{}
			err = d.Decode(list)
			if err != nil {
				fmt.Println(err)
				return
			}

			jsonOut, err := json.MarshalIndent(list, "", "\t")

			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonOut))
			return
		}

		var firstObject string
		if len(args) == 0 && stop {
			firstObject = id
		} else {
			firstObject = args[0]
		}
		var secondObject string

		if router {
			if len(args) == 2 {
				secondObject = args[1]
			} else {
				fmt.Println("Missing second router")
				return
			}
		}

		var req *http.Request
		var body []byte
		var err error

		if node {
			body, err = json.Marshal(&api.SniffNodeRequest{
				Node: firstObject,
				Id:   id,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/sniffNode", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if bridge {

			body, err = json.Marshal(&api.SniffBridgeRequest{
				Bridge: firstObject,
				Id:     id,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/sniffBridge", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if router {
			body, err = json.Marshal(&api.SniffRoutersRequest{
				Router1: firstObject,
				Router2: secondObject,
				Id:      id,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/sniffRouters", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if stop {
			body, err = json.Marshal(&api.StopSniffRequest{Id: firstObject})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/stopSniff", bytes.NewBuffer(body))
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
		if stop {
			resp = api.StopSniffResponse{}
		} else if node {
			resp = api.SniffNodeResponse{}
		} else if bridge {
			resp = api.SniffBridgeResponse{}
		} else if router {
			resp = api.SniffRoutersResponse{}
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
	rootCmd.AddCommand(sniffCmd)

	sniffCmd.Flags().BoolP("node", "n", false, "Sniffs node traffic")
	sniffCmd.Flags().BoolP("bridge", "b", false, "Sniffs bridge traffic")
	sniffCmd.Flags().BoolP("router", "r", false, "Sniffs traffic between routers")
	sniffCmd.Flags().BoolP("stop", "s", false, "Stops sniffing")
	sniffCmd.Flags().StringP("id", "i", "", "Component Id")
	sniffCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router", "stop")
}
