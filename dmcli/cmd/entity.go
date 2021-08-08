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
		Long:  `Generates entities for data model.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Generating content entities for " + util.AbsHomePath())
			return entity.Generate("entity")
		},
	}
)

func init() {
	rootCmd.AddCommand(entityCmd)
}
