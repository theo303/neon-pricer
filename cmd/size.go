package cmd

import (
	"fmt"
	"theo303/neon-pricer/conf"
	"theo303/neon-pricer/internal/usecases"

	"github.com/spf13/cobra"
)

// sizeCmd represents the sizeCmd command
var sizeCmd = &cobra.Command{
	Use:   "size",
	Short: "Calculate the global superficy of  a svg file.",
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
		bounds, err := usecases.GetBounds(formsGroups)
		if err != nil {
			panic(err)
		}
		for id, b := range bounds {
			width := b.Width()
			height := b.Height()
			fmt.Printf("%s: width=%.2fpx/%.2fmm height=%.2fpx/%.2fmm\n",
				id, width, width*1000/conf.Scale, height, height*1000/conf.Scale)
			// fmt.Printf("%s: %.2fpxˆ2, %.2fmmˆ2\n", id, size, size*1000/conf.Scale)
		}
	},
}

func init() {
	rootCmd.AddCommand(sizeCmd)

	sizeCmd.Flags().StringP("group", "g", "", "group id")
}
