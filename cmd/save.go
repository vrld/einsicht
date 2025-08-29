package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vrld/einsicht/internal"
)

var saveCmd = &cobra.Command{
	Use:     "save",
	Short:   "Save email parts as files",
	Aliases: []string{"s"},
	RunE:    saveTextCmd.RunE,
}

var saveTextCmd = &cobra.Command{
	Use:     "text",
	Short:   "Save text/plain to file",
	Aliases: []string{"t", "p"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.SaveToFileWithDialog("Save text body", theEmail.Subject+".txt", func(f *os.File) error {
			_, err := f.WriteString(theEmail.Text)
			return err
		})
	},
}

var saveHtmlCmd = &cobra.Command{
	Use:     "html",
	Short:   "Save cleaned text/html body to file",
	Aliases: []string{"h"},
	RunE: func(cmd *cobra.Command, args []string) error {
		html, err := internal.GetCleanedHTML(theEmail)
		if err != nil {
			return err
		}
		return internal.SaveToFileWithDialog("Save html body", theEmail.Subject+".html", func(f *os.File) error {
			_, err := f.WriteString(html)
			return err
		})
	},
}

var saveAttachmentCmd = &cobra.Command{
	Use:   "attachment <n>",
	Short: "Save attachment to file",
	Long: `Save attachment to file

Arguments:
  n    Which attachment to save (see einsicht list attachments)`,
	Aliases: []string{"a", "att"},
	Args:    internal.AttachmentSelectionArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		attachment, err := internal.AttachmentFromArg(args[0], theEmail.Attachments)
		if err != nil {
			return err
		}

		return internal.SaveToFileWithDialog("Save attachment", attachment.Filename, func(f *os.File) error {
			_, err := f.Write(attachment.Content)
			return err
		})
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.AddCommand(saveTextCmd, saveHtmlCmd, saveAttachmentCmd)
}
