package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vrld/einsicht/internal"
)

var printCmd = &cobra.Command{
	Use:     "print",
	Short:   "Print email part",
	Aliases: []string{"p"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(theEmail.Text)
	},
}

var printTextCmd = &cobra.Command{
	Use:     "text",
	Short:   "Print text/plain to stdout",
	Aliases: []string{"t", "p"},
	Run:     printCmd.Run,
}

var printHtmlCmd = &cobra.Command{
	Use:     "html",
	Short:   "Print cleaned text/html body to stdout",
	Aliases: []string{"h"},
	RunE: func(cmd *cobra.Command, args []string) error {
		html, err := internal.GetCleanedHTML(theEmail)
		if err != nil {
			return err
		}
		fmt.Print(html)
		return nil
	},
}

var printAttachmentCmd = &cobra.Command{
	Use:   "attachment <n>",
	Short: "Print attachment to stdout",
	Long: `Print attachment to stdout

Arguments:
  n    Which attachment to print (see einsicht list attachments)`,
	Aliases: []string{"a", "att"},
	Args:    internal.AttachmentSelectionArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		attachment, err := internal.AttachmentFromArg(args[0], theEmail.Attachments)
		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(attachment.Content)
		return err
	},
}

func init() {
	printCmd.AddCommand(printTextCmd, printHtmlCmd, printAttachmentCmd)
	rootCmd.AddCommand(printCmd)
}
