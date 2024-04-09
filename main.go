package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "rkl",
		Long:    "Rekall (rkl) is a CLI that helps you remember things.",
		Example: "rkl",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			viper.AutomaticEnv()
			viper.SetEnvPrefix("rkl")

			if _, err := os.Stat(viper.GetString(cmd.cfg)); errors.Is(err, os.ErrNotExist) {
				return errors.New(err.Error() + ": please run init to configure rkl\n")
			}
			return nil
		},
	}
	rootCmd.AddCommand(initialize())

	return rootCmd.ExecuteContext(context.Background())
}

func main() {
	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}
