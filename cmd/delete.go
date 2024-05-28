/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"goll/graph"
	"log"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		driver := graph.SelectDriver(dbselect)
		ctx := context.Background()

		defer driver.Driver.Close(ctx)
		var err error
		if err = graph.DBtableAllClear(ctx, driver.Driver, driver.DBname); err != nil {
			log.Fatal(err)
		}
	},
}

var dbselect string

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&dbselect, "select", "s", "", "input selected db")
	graphdbCmd.MarkFlagRequired("select")
}
