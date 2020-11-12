package main

import (
	"github.com/backup-operator/pkg/dumper"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var (
	config string
)

// NewBackupCommand return a full backup subcommand.
func NewBackupCommand() *cobra.Command {
	command := &cobra.Command{
		Use:          "backup",
		Short:        "backup a mysql",
		Example:      "Usage: mysqldumper backup -c conf/mydumper.ini.sample",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			arguments, err := dumper.ParseDumperConfig(config)
			if err != nil {
				klog.Fatalf("parse from config file failed")
			}

			dumper.Dumper(arguments)
		},
	}

	command.Flags().StringVarP(&config, "config", "c","", "config file")
	return command
}