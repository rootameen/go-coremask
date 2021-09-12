package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-coremask",
	Short: "Generate CPU Affinity Mask",
	Long: `A quick way to generate CPU Affinity Mask which can be 
used to pin processes to specific cpu cores
	
	Examples:
	  go-coremask generate --cores 44 --selection 2,4,6,8,10,12,14,16
	  go-coremask generate -c 8 -s 2,4 --hyperthreading=false
	  go-coremask detect`,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
