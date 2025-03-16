package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	api "github.com/David-Antunes/gone/api/Connect"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Short: "Connects network components together.",
	Use:   "connect [flags] [-l float | -w string | -d float | -c int | -j float ] {-n | -b | -r} <id> <id>",
	Example: `
	gone-cli connect -l 10 -w 100M -r router1 router2 
	
Connects router1 to router2 with a delay of 10 ms (5 ms for each direction) 
with a bandwidth limit of 100Mbits.`,
	Run: func(cmd *cobra.Command, args []string) {

		from := args[0]
		to := args[1]

		latency, _ := cmd.Flags().GetFloat64("latency")

		jitter, _ := cmd.Flags().GetFloat64("jitter")

		dropRate, _ := cmd.Flags().GetFloat64("dropRate")

		bandwidth, _ := cmd.Flags().GetString("bandwidth")

		weight, _ := cmd.Flags().GetInt("cost")
		
    propagate, _ :=  cmd.Flags().GetBool("propagate")

    propagate = !propagate

		node, _ := cmd.Flags().GetBool("node")

		bridge, _ := cmd.Flags().GetBool("bridge")

		router, _ := cmd.Flags().GetBool("router")

		bandwidthArr := []rune(bandwidth)
		multiplier := bandwidthArr[len(bandwidthArr)-1]
		bandwidthValue, err := strconv.Atoi(bandwidth[:len(bandwidthArr)-1])

		if err != nil {
			fmt.Println("Invalid value for bandwidth")
			return
		}
		if bandwidthValue <= 0 {
			fmt.Println("Invalid value for bandwidth")
			return
		}

		switch multiplier {
		case 'G':
			bandwidthValue *= 1000
			fallthrough
		case 'M':
			bandwidthValue *= 1000
			fallthrough
		case 'K':
			bandwidthValue *= 1000
		default:
			fmt.Println("Invalid bandwidth value")
			return
		}

		// if weight == -1 {
		// 	weight = bandwidthValue
		// }

		client := startClient()

		var req *http.Request
		var body []byte
		//		var err error

		if node {
			body, err = json.Marshal(&api.ConnectNodeToBridgeRequest{Node: from, Bridge: to, Latency: latency, Jitter: jitter, DropRate: dropRate, Bandwidth: bandwidthValue, Weight: weight})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/connectNodeToBridge", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if bridge {
			body, err = json.Marshal(&api.ConnectBridgeToRouterRequest{Bridge: from, Router: to, Latency: latency, Jitter: jitter, DropRate: dropRate, Bandwidth: bandwidthValue, Weight: weight})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/connectBridgeToRouter", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else if router {
      body, err = json.Marshal(&api.ConnectRouterToRouterRequest{From: from, To: to, Latency: latency, Jitter: jitter, DropRate: dropRate, Bandwidth: bandwidthValue, Weight: weight, Propagate: propagate})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/connectRouterToRouter", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Missing [ -n | -b | -r] flag")
			return
		}

		req.Header.Add("Content-Type", "application/json")

		jsonOutput(cmd, body, req)

		res, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		d := json.NewDecoder(res.Body)

		var resp any

		if node {
			resp = api.ConnectNodeToBridgeResponse{}
		} else if bridge {
			resp = api.ConnectBridgeToRouterResponse{}
		} else if router {
			resp = api.ConnectRouterToRouterResponse{}
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
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().Float64P("latency", "l", 0.0, "Configures link latency (ms)")
	connectCmd.Flags().Float64P("jitter", "j", 0.0, "Configures link jitter (ms)")
	connectCmd.Flags().Float64P("dropRate", "d", 0.0, "Configures link drop rate (Between 0.0 and 1.0)")
	connectCmd.Flags().StringP("bandwidth", "w", "10M", "Configures link bandwidth (Accepts 10K 10M or 10G)")
	connectCmd.Flags().IntP("cost", "c", 100, "Configures link weight")
	connectCmd.Flags().BoolP("propagate", "P", false, "Stop propagation when connecting two routers")

	connectCmd.Flags().BoolP("node", "n", false, "Connects a node to a bridge")
	connectCmd.Flags().BoolP("bridge", "b", false, "Connects a bridge to a router")
	connectCmd.Flags().BoolP("router", "r", false, "Connects two routers")
	connectCmd.MarkFlagsMutuallyExclusive("node", "bridge", "router")

}
