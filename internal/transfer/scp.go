package transfer

import (
	"fmt"
	"os"
	"path/filepath"

	"rioctl/internal/sshclient"
	"rioctl/internal/ui"
)

func PullFiles(client *sshclient.Client, remoteBase, localBase string, files []string) error {
	progressChan := make(chan int)

	// progress UI
	go ui.RunProgress(len(files), progressChan)

	for i, file := range files {
		remotePath := filepath.Join(remoteBase, file)
		localPath := filepath.Join(localBase, file)

		if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
			return err
		}

		// goph SCP download
		err := client.Raw().Download(remotePath, localPath)
		if err != nil {
			return fmt.Errorf("failed to download %s: %w", file, err)
		}

		progressChan <- i + 1
	}

	close(progressChan)
	return nil
}
