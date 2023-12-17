package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "neon-pricer",
	Short: "Calculate the price of a neon from a SVG file.",
	Long:  `Calculate the price of a neon from a SVG file.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
