/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// graphdbCmd represents the graphdb command
var graphdbCmd = &cobra.Command{
	Use:   "graphdb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var filepath string

func init() {
	rootCmd.AddCommand(graphdbCmd)
	// ファイルへのパスを指定する
	graphdbCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "input ast filepath")
	graphdbCmd.MarkFlagRequired("filepath")
}
