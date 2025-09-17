package cmd

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vrld/einsicht/internal"
	"github.com/vrld/einsicht/internal/ui"
	"golang.org/x/term"
)

// package local variable! this is shared between all commands and guaranteed to be non-null on Run or RunE
var theEmail *internal.Email

var rootCmd = &cobra.Command{
	Use:               "einsicht",
	Short:             "The mail reader and multitool",
	Long:              "The mail reader and multitool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetEnvPrefix("EINSICHT")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
		viper.AutomaticEnv()
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		return loadEmailFromFlags(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.Run(theEmail)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("file", "f", "-", "read email from `PATH`; `-` means stdin")
	rootCmd.Flags().StringP("body", "b", "plain", "Preferred body to display: plain or html")
	rootCmd.Flags().StringP("command", "c", "xdg-open", "Open files with this command")
}

func loadEmailFromFlags(cmd *cobra.Command) error {
	emailPath := cmd.Flag("file").Value.String()
	if emailPath == "-" {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return setTheEmail(bufio.NewReader(os.Stdin))
		}
	}

	file, err := os.Open(emailPath)
	if err == nil {
		err = setTheEmail(file)
		file.Close()
	}
	return err
}

func setTheEmail(reader io.Reader) error {
	mail, err := internal.ReadEmail(reader)
	if err != nil {
		return err
	}
	theEmail = mail
	return nil
}
