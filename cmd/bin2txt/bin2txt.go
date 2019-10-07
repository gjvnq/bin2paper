package main

import (
	"fmt"
	"os"

	"../.."
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bin2txt [input] [output]",
	Short: "Converts a binary file to a printable .txt",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		output := ""
		if len(args) >= 2 {
			output = args[1]
		}

		err := bin2paper.TXTFromFile(input, output)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
