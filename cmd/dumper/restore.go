package main

import (
	"github.com/backup-operator/pkg/dumper"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func NewRestoreCommand() *cobra.Command {
	command := &cobra.Command{
		Use:          "restore",
		Short:        "restore a mysql",
		Example: 	  "Usage: mysqldumper restore -c conf/mydumper.ini.sample",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			arguments, err := dumper.ParseDumperConfig(config)
			if err != nil {
				klog.Fatalf("parse from config file failed")
			}
			dumper.Loader(arguments)
		},
	}

	command.Flags().StringVarP(&config, "config", "c","", "config file")
	return command
}