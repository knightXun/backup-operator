package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/backup-operator/pkg/to-crdgen/cmd"
	"github.com/spf13/pflag"
)

func main() {
	flags := pflag.NewFlagSet("to-crdgen", pflag.ExitOnError)
	flag.CommandLine.Parse([]string{})
	pflag.CommandLine = flags

	command := cmd.NewToCrdGenRootCmd()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}