package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate coremask based on user input",
	Long: `Generate coremask based on user input. This requires: 
- Number of Cores in the target system
- Whether Hyper Threading is enabled or not
- The selection of cores where the affinity should be set`,
	Run: func(cmd *cobra.Command, args []string) {

		numCores, _ := cmd.Flags().GetInt("cores")
		coreSelection, _ := cmd.Flags().GetString("selection")
		cores := make([]string, numCores)
		selection := strings.Split(coreSelection, ",")
		for core := range cores {
			if coreMatch(selection, strconv.Itoa(core)) {
				cores[core] = "1"
			} else {
				cores[core] = "0"
			}
		}
		coremaskBinary := strings.Join(reverseSlice(cores), "")

		htStatus, _ := cmd.Flags().GetBool("hyperthreading")

		// On systems with Hyperthreading, the coremask is adjusted to reflect the siblings of each CPU
		if htStatus {
			coremaskBinary += coremaskBinary
		}

		coremask, err := strconv.ParseInt(string(coremaskBinary), 2, 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("* Total number of CPUs is: %d\n", len(coremaskBinary))
		fmt.Printf("* Selected cores for affinity are: %q\n", coreSelection)
		fmt.Printf("* Here's the affinity coremask: 0x%x\n", coremask)

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().BoolP("hyperthreading", "t", true, "HyperThreading / SMT Enabled")
	generateCmd.Flags().IntP("cores", "c", 8, "Number of cores in system")
	generateCmd.Flags().StringP("selection", "s", "", "Selection of cores for for affinity (comma-separated, e.g. `2,4,6,8`)")
	generateCmd.MarkFlagRequired("cores")
	generateCmd.MarkFlagRequired("selection")
}

func coreMatch(cores []string, selection string) bool {
	sort.Strings(cores)
	for _, core := range cores {
		if core == selection {
			return true
		}
	}
	return false
}

func reverseSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
