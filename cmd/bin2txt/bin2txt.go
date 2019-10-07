package main

import (
  "fmt"
  "os"
  "github.com/spf13/cobra"
  "../.."
)

var rootCmd = &cobra.Command{
  Use:   "bin2txt",
  Short: "Converts a binary file to a printable .txt",
  Run: func(cmd *cobra.Command, args []string) {
  	err := bin2paper.TXTFromFile(args[0])
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