package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/David-Antunes/gone/api/Connect"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Args:  cobra.ExactArgs(3),
	Short: "Connects to network components together.",
	Long:  `Connects to network components together in the emulation with the properties defined`,
	Run: func(cmd *cobra.Command, args []string) {

		netType := args[0]
		from := args[1]
		to := args[2]

		latency, _ := cmd.Flags().GetInt("latency")

		jitter, _ := cmd.Flags().GetFloat64("jitter")

		dropRate, _ := cmd.Flags().GetFloat64("dropRate")

		bandwidth, _ := cmd.Flags().GetString("bandwidth")

		weight, _ := cmd.Flags().GetInt("weight")

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

		if weight == -1 {
			weight = bandwidthValue
		}

		client := startClient()

		var req *http.Request
		switch netType {
		case "node":
			body, err := json.Marshal(&api.ConnectNodeToBridgeRequest{Node: from, Bridge: to, Latency: latency, Jitter: jitter, DropRate: dropRate, Bandwidth: bandwidthValue, Weight: weight})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/connectNodeToBridge", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			jsonOutput(cmd, body, req)
		case "bridge":
			body, err := json.Marshal(&api.ConnectBridgeToRouterRequest{Bridge: from, Router: to, Latency: latency, Jitter: jitter, DropRate: dropRate, Bandwidth: bandwidthValue, Weight: weight})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/connectBridgeToRouter", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			jsonOutput(cmd, body, req)
		case "router":
			body, err := json.Marshal(&api.ConnectRouterToRouterRequest{From: from, To: to, Latency: latency, Jitter: jitter, DropRate: dropRate, Bandwidth: bandwidthValue, Weight: weight})
			if err != nil {
				panic(err)
			}
			req, err = http.NewRequest("POST", URL+"/connectRouterToRouter", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			jsonOutput(cmd, body, req)
		default:
			panic("Something really went wrong!")
		}

		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		d := json.NewDecoder(res.Body)

		switch netType {
		case "node":
			resp := &api.ConnectNodeToBridgeResponse{}
			err = d.Decode(&resp)

			if err != nil {
				panic(err)
			}
			prettyPrint(resp)

		case "bridge":
			resp := &api.ConnectBridgeToRouterResponse{}
			err = d.Decode(&resp)

			if err != nil {
				panic(err)
			}
			prettyPrint(resp)
		case "router":
			resp := &api.ConnectRouterToRouterResponse{}
			err = d.Decode(&resp)

			if err != nil {
				panic(err)
			}
			prettyPrint(resp)
		default:
			panic("Something really went wrong!")
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().IntP("latency", "l", 0, "Configures link latency")
	connectCmd.Flags().Float64P("jitter", "j", 0.0, "Configures link jitter")
	connectCmd.Flags().StringP("bandwidth", "b", "1M", "Configures link bandwidth")
	connectCmd.Flags().Float64P("dropRate", "d", 0.0, "Configures link drop rate")
	connectCmd.Flags().IntP("weight", "w", -1, "Configures link weight")
}
