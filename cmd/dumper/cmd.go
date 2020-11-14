package main

import (
	"github.com/backup-operator/pkg/dumper"
	"k8s.io/klog"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	config string
)

func main() {
	rootCmd := &cobra.Command{
		Use:              "mysqldumper",
		Short:            "mysqldumper dumper/restore tools.",
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	rootCmd.AddCommand(
		NewBackupCommand(),
		NewRestoreCommand(),
	)

	rootCmd.SetOut(os.Stdout)

	rootCmd.SetArgs(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		klog.Error("mysqldumper failed", zap.Error(err))
		os.Exit(1)
	}
}

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

// NewBackupCommand return a full mydumper subcommand.
func NewBackupCommand() *cobra.Command {
	command := &cobra.Command{
		Use:          "mydumper",
		Short:        "mydumper a mysql",
		Example:      "Usage: mysqldumper mydumper -c conf/mydumper.ini.sample",
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