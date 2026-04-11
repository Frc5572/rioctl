package transfer

import (
	"os"
	"path/filepath"

	"rioctl/internal/sshclient"
	"rioctl/internal/ui"
)

func PullFiles(client *sshclient.Client, remoteBase, localBase string, files []string) error {
	progressChan := make(chan int)

	go func() {
		for i, file := range files {
			remotePath := filepath.Join(remoteBase, file)
			localPath := filepath.Join(localBase, file)

			if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
				// continue
			}

			// goph SCP download
			err := client.Raw().Download(remotePath, localPath)
			if err != nil {
				// fmt.Printf("failed to download %s: %w", file, err)
				// continue
			}

			progressChan <- i + 1
		}
	}()
	// progress UI
	ui.RunProgress(len(files), progressChan)

	close(progressChan)
	return nil
}
