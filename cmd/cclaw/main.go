package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/actuagent/cclaw/pkg/app"
	"github.com/actuagent/cclaw/pkg/config"
	"github.com/spf13/cobra"
)

var (
	// Version is the build version, overridden via -ldflags at build time.
	Version = "dev"
	cfgFile string
)

func main() {
	app.InitLogger()

	rootCmd := &cobra.Command{
		Use:   "cclaw",
		Short: "AI agent with tools, memory, and channels",
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.cclaw/config.yaml)")

	rootCmd.AddCommand(
		newGatewayCmd(),
		newAgentCmd(),
		newOnboardCmd(),
		newStatusCmd(),
		newVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func mustLoadConfig() *config.Config {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	return cfg
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cclaw %s\n", Version)
		},
	}
}
