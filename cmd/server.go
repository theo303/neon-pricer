package cmd

import (
	"theo303/neon-pricer/conf"
	"theo303/neon-pricer/internal/http"

	"github.com/spf13/cobra"
)

// serverCmd represents the serverCmd command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server.",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := conf.Load()
		if err != nil {
			panic(err)
		}

		api := http.NewAPI(conf, 8080)

		if err = api.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
