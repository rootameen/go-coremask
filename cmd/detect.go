package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type cpu struct {
	sockets int
	cpus    int
	numas   int
	ht      bool
}

// detectCmd represents the detect command
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect the current core distribution on a system",
	Long: `Detect the current core distribution on a system. This makes 
it easier to get the inputs for "go-coremask generate" command. 
Ideally, detect should be used on the target system where affinity is to
be set, or an exact replica to avoid generating an incorrect coremask.`,
	Run: func(cmd *cobra.Command, args []string) {

		// accept detect command only on Linux systems
		os := runtime.GOOS
		switch os {
		case "linux":
			fmt.Println("= System Information =")
			cpu := cpuInfo()
			fmt.Printf("Number of Physical CPU Sockets: %d\n", cpu.sockets)
			fmt.Printf("Number of CPU Cores (Per Socket): %d\n", cpu.cpus)
			if cpu.ht {
				fmt.Printf("** HyperThreading Enabled! Total Number of CPUs is %d\n", cpu.cpus*cpu.sockets)
			} else {
				fmt.Printf("** HyperThreading Disabled! Total Number of CPUs is %d\n", cpu.cpus)
			}
			fmt.Printf("Number of NUMA Nodes: %d\n", cpu.numas)
		default:
			fmt.Println("`go-cmdline detect` is currently supported on Linux")
		}

	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}

func infoParse(s string) int {
	count := strings.Split(strings.TrimSpace(s), "-")
	maxVal, _ := strconv.Atoi(count[1])
	return (maxVal + 1)
}

func readFile(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func cpuInfo() cpu {
	c := cpu{}

	// Check Hyperthreading status
	if strings.TrimSpace(readFile("/sys/devices/system/cpu/smt/active")) == "1" {
		c.ht = true
	} else {
		c.ht = false
	}

	// Check number of active NUMA nodes
	if strings.TrimSpace(readFile("/sys/devices/system/node/online")) == "0" {
		c.numas = 1
	} else {
		c.numas = infoParse(readFile("/sys/devices/system/node/online"))
	}

	// Check number of Physical CPU Sockets
	cmd := exec.Command("sh", "-c", "grep -i 'physical id' /proc/cpuinfo | sort -u | wc -l")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	sockets := strings.TrimSpace(string(stdout))
	c.sockets, _ = strconv.Atoi(sockets)

	// Check Total number of of CPUs
	cmd = exec.Command("sh", "-c", "grep -i 'cpu cores' /proc/cpuinfo | sort -u | rev | cut -d' ' -f1 | rev")
	stdout, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	cpus := strings.TrimSpace(string(stdout))
	cpusVal, _ := strconv.Atoi(cpus)
	c.cpus = cpusVal * c.sockets

	return c
}
