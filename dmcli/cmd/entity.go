package cmd

import (
	"fmt"

	"github.com/digimakergo/digimaker/codegen/entity"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/cobra"
)

var (
	entityCmd = &cobra.Command{
		Use:   "entity",
		Short: "Generates entities",
		Long:  `Generates entities for data model. usage: dmcli entity <folder>(if not provided, 'entity' will be used)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := LoadPlugin(cmd)
			if err != nil {
				return err
			}
			folder := "entity"
			if len(args) > 0 {
				folder = args[0]
			}
			fmt.Println("Generating content entities for " + util.AbsHomePath())
			return entity.Generate(folder)
		},
	}
)

func init() {
	rootCmd.AddCommand(entityCmd)
}
