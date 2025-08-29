package cmd

import (
	"bufio"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/vrld/einsicht/internal"
	"github.com/vrld/einsicht/internal/ui"
	"golang.org/x/term"
)

// package local variable! this is shared between all commands and guaranteed to be non-null on Run or RunE
var theEmail *internal.Email

var rootCmd = &cobra.Command{
	Use:               "einsicht",
	Short:             "The mail reader and multitool",
	Long:              "Read your mail with style; process it with ease",
	PersistentPreRunE: loadEmailFromFlags,
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
}

func loadEmailFromFlags(cmd *cobra.Command, args []string) error {
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
