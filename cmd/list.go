package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrld/einsicht/internal"
)

var asJson bool  // TODO: implement this

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List email parts",
	Aliases: []string{"l"},
	Run:     listPartsCmd.Run,
}

var listPartsCmd = &cobra.Command{
	Use:     "parts",
	Short:   "List email parts",
	Aliases: []string{"p"},
	Run: func(cmd *cobra.Command, args []string) {
		listBodiesCmd.Run(cmd, args)
		listAttachmentsCmd.Run(cmd, args)
	},
}

var listBodiesCmd = &cobra.Command{
	Use:     "bodies",
	Short:   "List email bodies",
	Aliases: []string{"b"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(theEmail.Text) > 0 {
			fmt.Println("- text/plain ", internal.HumanReadableFileSize(len(theEmail.Text)))
		}
		if len(theEmail.HTML) > 0 {
			fmt.Println("- text/html ", internal.HumanReadableFileSize(len(theEmail.HTML)))
		}
	},
}

var listAttachmentsCmd = &cobra.Command{
	Use:     "attachments",
	Short:   "List email attachments",
	Aliases: []string{"a", "att"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(theEmail.Attachments) == 0 {
			fmt.Println("- No attachments")
			return
		}

		fmt.Println("- Attachments:")
		for i, a := range theEmail.Attachments {
			fmt.Printf("  %d. %s (%s) %s\n", i+1, a.Filename, a.ContentType, internal.HumanReadableFileSize(len(a.Content)))
		}
	},
}

func init() {
	listCmd.AddCommand(listPartsCmd, listBodiesCmd, listAttachmentsCmd)
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolVarP(&asJson, "json", "j", false, "output as json")
}
