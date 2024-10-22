/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/David-Antunes/gone/api/Operations"
	"github.com/spf13/cobra"
	"net/http"
)

// disconnectCmd represents the disconnect command
var sniffCmd = &cobra.Command{
	Use:   "sniff",
	Short: "Sniffs traffic from Link",
	Long:  ``,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		client := startClient()

		if len(args) == 0 {
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
		}

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

		if node {
			if !stop {
				body, err := json.Marshal(&api.SniffNodeRequest{Name: firstObject})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/sniffNode", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			} else {
				body, err := json.Marshal(&api.StopSniffRequest{Id: firstObject})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/stopSniffNode", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			}
		} else if bridge {
			if !stop {

				body, err := json.Marshal(&api.SniffBridgeRequest{Name: firstObject})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/sniffBridge", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			} else {

				body, err := json.Marshal(&api.StopSniffRequest{Id: firstObject})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/stopSniffBridge", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			}
		} else if router {
			if !stop {

				body, err := json.Marshal(&api.SniffRoutersRequest{Router1: firstObject, Router2: secondObject})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/sniffRouters", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			} else {
				body, err := json.Marshal(&api.StopSniffRequest{Id: firstObject})
				if err != nil {
					panic(err)
				}
				req, err = http.NewRequest("POST", URL+"/stopSniffRouters", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println("Missing -n/-b/-r flag")
			return
		}

		req.Header.Add("Content-Type", "application/json")

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
	sniffCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router")
}
