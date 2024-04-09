package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nandanjp/rkl/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	HistoryFile string
}

func ToFile(path string, config *Config) error {
	return nil
}

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:     "initialize",
		Short:   "init the rkl cfg.",
		Long:    "init provision the rkl configuration file.",
		Example: "rkl init",
		Aliases: []string{"i", "init"},
		// used to overwrite/skip the parent commands persistentPreRunE func
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind Cobra flags with viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Environment variables are expected to be ALL CAPS
			viper.AutomaticEnv()
			viper.SetEnvPrefix("rkl")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			if err := tea.Run(tui.NewInitPrompt(viper.GetString(cfgPath), homeDir)).Start(); err != nil {
				return err
			}
			return nil
		},
	}
	return init
}
