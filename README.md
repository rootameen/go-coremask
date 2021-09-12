# go-coremask

A quick way to generate CPU Affinity Mask which can be used to pin processes to specific cpu cores

## Affinity Mask?

CPU Affinity is represented as a bitmask, in which each CPU core is representing by a bit in the mask. If the bit value is set to 1, then the process/thread is set to run on that core. By default, the value of all bits is set to 1, meaning all processes can run on any core in the system

## How it works?
Currently, two modes are supported:
* `generate` to generate mask based on user inputs
* `detect` as a helper to collect info as inputs for `generate` command


## Examples:

```
// generate coremask for a 44 cores system with hyperthreading on (default)
$ go-coremask generate --cores 44 --selection 2,4,6,8,10,12,14,16
* Total number of CPUs is: 88
* Selected cores for affinity are: "2,4,6,8,10,12,14,16"
* Here's the affinity coremask: 0x1555400000015554

// generate coremask for a system without hyperthreading enabled
$ go-coremask generate -c 8 -s 2,4 --hyperthreading=false
* Total number of CPUs is: 8
* Selected cores for affinity are: "2,4"
* Here's the affinity coremask: 0x14

// detect system CPU info
$ go-coremask detect
= System Information =
Number of Physical CPU Sockets: 2
Number of CPU Cores (Per Socket): 16
** HyperThreading Enabled! Total Number of CPUs is 32
Number of NUMA Nodes: 2
```
