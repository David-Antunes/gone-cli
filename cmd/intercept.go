/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/David-Antunes/gone/api/Operations"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"net/http"
)

// disconnectCmd represents the disconnect command
var interceptCmd = &cobra.Command{
	Use:   "intercept",
	Short: "Intercepts traffic from Link",
	Long:  ``,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		client := startClient()

		if len(args) == 0 {
			req, err := http.NewRequest("GET", URL+"/listIntercepts", nil)
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
			list := &api.ListInterceptsResponse{}
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

		firstObject := args[0]
		var secondObject string
		node, _ := cmd.Flags().GetBool("node")

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		stop, _ := cmd.Flags().GetBool("stop")

		tx, _ := cmd.Flags().GetBool("receive")

		id, _ := cmd.Flags().GetString("id")

		if id == "" {
			u, err := uuid.NewUUID()
			if err != nil {
				panic(err)
			}
			id = u.String()
		}

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
			body, err = json.Marshal(&api.InterceptNodeRequest{
				Node:      firstObject,
				Id:        id,
				Direction: tx,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/interceptNode", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if bridge {
			body, err = json.Marshal(&api.InterceptBridgeRequest{
				Bridge:    firstObject,
				Id:        id,
				Direction: tx,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/interceptBridge", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if router {

			body, err = json.Marshal(&api.InterceptRoutersRequest{
				Router1:   firstObject,
				Router2:   secondObject,
				Id:        id,
				Direction: tx,
			})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/interceptRouter", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if stop {

			body, err = json.Marshal(&api.StopInterceptRequest{Id: firstObject})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/stopIntercept", bytes.NewBuffer(body))
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
			resp = api.StopInterceptResponse{}
		} else if node {
			resp = api.InterceptNodeResponse{}
		} else if bridge {
			resp = api.InterceptBridgeResponse{}
		} else if router {
			resp = api.InterceptRoutersResponse{}
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
	rootCmd.AddCommand(interceptCmd)

	interceptCmd.Flags().BoolP("node", "n", false, "intercepts node traffic")
	interceptCmd.Flags().BoolP("bridge", "b", false, "intercepts bridge traffic")
	interceptCmd.Flags().BoolP("router", "r", false, "intercepts traffic between routers")
	interceptCmd.Flags().BoolP("stop", "s", false, "Stops sniffing")
	interceptCmd.Flags().BoolP("receive", "x", true, "intercepts receive traffic")
	interceptCmd.Flags().StringP("id", "i", "", "Component Id")
	interceptCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router", "stop")
}
