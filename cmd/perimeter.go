/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"theo303/neon-pricer/internal/usecases"

	"github.com/spf13/cobra"
)

// perimeterCmd represents the perimeter command
var perimeterCmd = &cobra.Command{
	Use:   "perimeter",
	Short: "Calculate the total perimeter of all forms in a svg file.",
	Long: `Calculate the total perimeter of all forms in a svg file.
	
	For now, only rectangles and circles are supported`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		perimeter, err := usecases.GetPerimeter(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Println(perimeter)
	},
}

func init() {
	rootCmd.AddCommand(perimeterCmd)
}
