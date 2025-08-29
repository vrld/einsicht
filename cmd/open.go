package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vrld/einsicht/internal"
)

var openCmd = &cobra.Command{
	Use:     "open",
	Short:   "Open email parts",
	Aliases: []string{"o"},
	RunE:    openTextCmd.RunE,
}

var openTextCmd = &cobra.Command{
	Use:     "text",
	Short:   "Open text/plain body",
	Aliases: []string{"t", "p"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.OpenStringWithCommand(openCommand, "txt", func(file *os.File) error {
			_, err := file.WriteString(theEmail.Text)
			return err
		})
	},
}

var openHtmlCmd = &cobra.Command{
	Use:     "html",
	Short:   "Open cleaned text/html body",
	Aliases: []string{"h"},
	RunE: func(cmd *cobra.Command, args []string) error {
		html, err := internal.GetCleanedHTML(theEmail)
		if err != nil {
			return err
		}
		return internal.OpenStringWithCommand(openCommand, "html", func(file *os.File) error {
			_, err := file.WriteString(html)
			return err
		})
	},
}

var openAttachmentCmd = &cobra.Command{
	Use:   "attachment <n>",
	Short: "Open attachment",
	Long: `Open attachment

Arguments:
  n    Which attachment to open (see einsicht list attachments)`,
	Aliases: []string{"a", "att"},
	Args:    internal.AttachmentSelectionArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		attachment, err := internal.AttachmentFromArg(args[0], theEmail.Attachments)
		if err != nil {
			return err
		}

		return internal.OpenStringWithCommand(openCommand, attachment.Filename, func(file *os.File) error {
			_, err = file.Write(attachment.Content)
			return err
		})
	},
}

var openCommand string

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.PersistentFlags().StringVarP(&openCommand, "command", "c", "xdg-open", "Open with this command")

	openCmd.AddCommand(openTextCmd, openHtmlCmd, openAttachmentCmd)
}
