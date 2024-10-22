/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restarts gone in all participating machines",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if resp, err := http.Get(AGENT_URL + "/restart"); err != nil {
			fmt.Println(resp.Body)
		} else {
			fmt.Println(err)
		}

	},
}

func init() {

}
