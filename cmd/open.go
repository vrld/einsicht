package cmd

import (
	"os"
	"os/exec"
	"time"

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
		return openContent("txt", func(file *os.File) error {
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
		return openContent("html", func(file *os.File) error {
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

		return openContent(attachment.Filename, func(file *os.File) error {
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

func openContent(suffix string, writer func(*os.File) error) error {
	file, err := os.CreateTemp("", "*."+suffix)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	err = writer(file)
	file.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command(openCommand, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	time.Sleep(time.Second * 1)  // give the process some time to read the file before deleting it
	return err
}
