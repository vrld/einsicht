package internal

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rymdport/portal/filechooser"
)

type WriteContentFunc func(*os.File) error

func OpenStringWithCommand(
	openCommand string,
	fileEnding string,
	writeContent WriteContentFunc,
) error {
	file, err := os.CreateTemp("", "*."+fileEnding)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	err = writeContent(file)
	file.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command(openCommand, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		time.Sleep(time.Second * 1) // firefox returns immediately, but needs some time to open the file
	}
	return err
}

func SaveToFileWithDialog(
	dialogTitle,
	suggestedFilename string,
	writeContext WriteContentFunc,
) error {
	options := filechooser.SaveFileOptions{CurrentName: suggestedFilename}
	files, err := filechooser.SaveFile("einsicht", dialogTitle, &options)
	if err != nil {
		return err
	}

	for _, filename := range files {
		filename := strings.TrimPrefix(filename, "file://")
		if file, err := os.Create(filename); err != nil {
			return err
		} else if err := writeContext(file); err != nil {
			return err
		}
	}

	return nil
}
