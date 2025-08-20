package cmd

import (
	"os"
	"strings"

	"github.com/rymdport/portal/filechooser"
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
		return savePart([]byte(theEmail.Text), "Save text body", theEmail.Subject+".txt")
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
		return savePart([]byte(html), "Save html body", theEmail.Subject+".html")
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

		return savePart(attachment.Content, "Save attachment", attachment.Filename)
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.AddCommand(saveTextCmd, saveHtmlCmd, saveAttachmentCmd)
}

func savePart(content []byte, title, suggestedFilename string) error {
	options := filechooser.SaveFileOptions{CurrentName: suggestedFilename}
	files, err := filechooser.SaveFile("einsicht", title, &options)
	if err != nil {
		return err
	}

	for _, filename := range files {
		filename := strings.TrimPrefix(filename, "file://")
		if err = os.WriteFile(filename, content, os.ModePerm); err != nil {
			return err
		}
	}

	return nil

}
