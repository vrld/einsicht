package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func AttachmentByCID(email *Email, cid string) *Attachment {
	cid = strings.TrimPrefix(cid, "cid:")
	for _, attachment := range email.Attachments {
		if attachment.Cid == cid {
			return attachment
		}
	}
	return nil
}

func AttachmentSelectionArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	_, err := strconv.Atoi(args[0])
	return err
}

func AttachmentFromArg(arg string, attachments []*Attachment) (*Attachment, error) {
	index, err := strconv.Atoi(arg)
	if err != nil {
		return nil, err
	}

	if index <= 0 || index > len(attachments) {
		return nil, fmt.Errorf("No attachment with index %d. Valid range: [1:%d]", index, len(attachments))
	}

	return attachments[index-1], nil
}
