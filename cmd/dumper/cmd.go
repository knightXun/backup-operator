package main

import (
	"k8s.io/klog"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	rootCmd := &cobra.Command{
		Use:              "mysqldumper",
		Short:            "mysqldumper dumper/restore tool.",
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
