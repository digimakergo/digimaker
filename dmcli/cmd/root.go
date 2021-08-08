package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "dmcli",
		Short: "Cli tools for digimaker cmf",
		Long: `dmtool is cli tool for digimaker Content Management Framework. 
It helps to generate database entities, database schema sql, etc`,
	}
)

func Execute() {
	rootCmd.Execute()
}
