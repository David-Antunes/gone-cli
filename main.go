package main

import (
	"github.com/David-Antunes/gone-cli/cmd"
	"github.com/spf13/viper"
)

func main() {

	viper.SetDefault("URL", "http://localhost:3000")
	viper.SetDefault("AGENT_URL", "http://localhost:3300")
	viper.AutomaticEnv()
	cmd.URL = viper.GetString("URL")
	cmd.AGENT_URL = viper.GetString("AGENT_URL")
	cmd.Execute()
}
