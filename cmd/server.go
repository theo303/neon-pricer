package cmd

import (
	"theo303/neon-pricer/configuration"
	"theo303/neon-pricer/internal/api"

	"github.com/spf13/cobra"
)

// serverCmd represents the serverCmd command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server.",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := configuration.Load()
		if err != nil {
			panic(err)
		}

		err = api.Run(api.Conf{
			Configuration: conf,
			Port:          8080,
		})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
