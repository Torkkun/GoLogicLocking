/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"goll/graph"
	"goll/parser"
	"log"

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
		parseresult := parser.NewParse(filepath)
		// Connectionの設定は後で考える
		dbname := "neo4j"

		//driver := graph.NewDriver()
		driver := graph.SelectDriver("origin")
		ctx := context.Background()

		defer driver.Driver.Close(ctx)
		var err error
		for _, io := range parseresult.Declarations.IOPorts {
			neoio := graph.IONode{
				Type: string(io.Type),
				Name: io.Name,
			}
			if err = neoio.CreateInOutNode(ctx, driver.Driver, dbname); err != nil {
				log.Fatalln(err)
			}
		}
		for _, wire := range parseresult.Declarations.Wires {
			neowire := graph.WireNode{
				Name: wire.Name,
			}
			if err = neowire.CreateWireNode(ctx, driver.Driver, dbname); err != nil {
				log.Fatalln(err)
			}
		}
		for _, logicgate := range parseresult.LogicGates {
			gate := graph.LogicGateNode{
				GateType: string(logicgate.GateType),
				At:       logicgate.At,
			}
			if err = gate.CreateLogicGateNode(ctx, driver.Driver, dbname); err != nil {
				log.Fatalln(err)
			}
		}
		for at, relation := range parseresult.Nodes {
			// i1 <- lg, i2 <- lg, lg <- out
			// i1 <- lg
			if err := graph.LGtoIN(ctx, driver.Driver, dbname, relation.In1, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
				log.Fatalln(err)
			}

			// i2 <- lg
			if err := graph.LGtoIN(ctx, driver.Driver, dbname, relation.In2, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
				log.Fatalln(err)
			}

			// lg <- out
			if err := graph.OUTtoLG(ctx, driver.Driver, dbname, relation.Out, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
				log.Fatalln(err)
			}
		}
	},
}

var filepath string

func init() {
	rootCmd.AddCommand(graphdbCmd)
	// ファイルへのパスを指定する
	graphdbCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "input ast filepath")
	graphdbCmd.MarkFlagRequired("filepath")
}
