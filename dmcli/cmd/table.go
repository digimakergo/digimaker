package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/digimakergo/digimaker/codegen/table"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/spf13/cobra"
)

var (
	tableCmd = &cobra.Command{
		Use:   "db-table",
		Short: "Generates database tables",
		Long:  `Generates database table, args: <contenttype 1>,<contenttype 2>. With empty arg will generate all tables. note: only mysql is supported for now`,
		Run: func(cmd *cobra.Command, args []string) {
			var contenttypes []string
			if len(args) == 0 {
				contenttypeList := definition.GetDefinitionList()["default"]
				for identifier, _ := range contenttypeList {
					contenttypes = append(contenttypes, identifier)
				}
				sort.Strings(contenttypes)
			} else {
				contenttypes = args
			}

			fmt.Println("Generating table for " + strings.Join(contenttypes, ","))
			fmt.Println("----------")
			for _, contenttype := range contenttypes {
				err := table.GenerateTable(contenttype)
				if err != nil {
					fmt.Println("Fail to generate: " + err.Error())
				}
				fmt.Println("")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(tableCmd)
}
