package cmd

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/cmd/gen"
	"github.com/MarchGe/go-admin-server/cmd/server"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "go-admin-server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("error: requires at least one arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// nothing to do, only aim to throw the error above
	},
	Version: config.Version,
}

func init() {
	rootCmd.AddCommand(server.Server)
	rootCmd.AddCommand(gen.Gen)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
