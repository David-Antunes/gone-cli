package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

// shutdownCmd represents the shutdown command
var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdowns gone and gone-agent in all participating machines",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if resp, err := http.Get(AGENT_URL + "/shutdown"); err != nil {
			fmt.Println(resp.Body)
		} else {
			fmt.Println(err)
		}
	},
}

func init() {

}
