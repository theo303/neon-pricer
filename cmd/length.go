package cmd

import (
	"fmt"
	"theo303/neon-pricer/conf"
	"theo303/neon-pricer/internal/usecases"

	"github.com/spf13/cobra"
)

// lengthCmd represents the lengthCmd command
var lengthCmd = &cobra.Command{
	Use:   "length",
	Short: "Calculate the total length of all forms in a svg file.",
	Long: `Calculate the total length of all forms in a svg file.
Each groups of forms will be measured independantly and then summed together.
	
Rectangles, circles and paths are supported.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)

		conf, err := conf.Load()
		if err != nil {
			panic(err)
		}

		groupID, err := cmd.Flags().GetString("group")
		if err != nil {
			panic(err)
		}

		formsGroups, err := usecases.ParseSVGFile(args[0], groupID)
		if err != nil {
			panic(err)
		}
		lengths, err := usecases.GetLengths(formsGroups)
		if err != nil {
			panic(err)
		}
		var totalLength float64
		for id, length := range lengths {
			fmt.Printf("%s: %.2fpx, %.2fmm\n", id, length, length*1000/conf.Scale)
			totalLength += length
		}
		fmt.Printf("total: %.2f\n", totalLength)
	},
}

func init() {
	rootCmd.AddCommand(lengthCmd)

	lengthCmd.Flags().StringP("group", "g", "", "group id")
}
