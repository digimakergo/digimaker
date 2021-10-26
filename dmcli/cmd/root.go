package cmd

import (
	"fmt"
	"plugin"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dmcli",
		Short: "Cli tools for digimaker cmf",
		Long: `dmcli is cli tool for digimaker Content Management Framework. 
It helps to generate database entities, database schema sql, etc`,
	}
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("plugin", "", "Load .so for plugin fieldtypes and others")
}

func LoadPlugin(cmd *cobra.Command) error {
	//load plugin
	pPath := cmd.Root().PersistentFlags().Lookup("plugin").Value.String()
	if pPath != "" {
		_, err := plugin.Open(pPath)
		if err != nil {
			return fmt.Errorf("Error when loading plugin: %v", err)
		}
	}
	return nil
}
