package transfer

import (
	"fmt"
	"os"
	"path/filepath"

	"rioctl/internal/sshclient"
	"rioctl/internal/ui"
	"rioctl/internal/utils"
)

func PullFiles(client *sshclient.Client, remoteBase, localBase string, files []string) error {
	finishedChan := make(chan int)
	currentChan := make(chan utils.FileUpload)

	go func() {
		for i, file := range files {
			remotePath := filepath.Join(remoteBase, file)
			localPath := filepath.Join(localBase, file)

			if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
				// continue
			}

			err := os.Remove(localPath)
			if err != nil {
				//
			}
			totalSize, err := client.FileSize(remotePath)
			if err != nil {
				//
			}
			file := utils.NewFile(totalSize, 0, file)
			currentChan <- file

			go func() {
				for {
					fi, err := os.Stat(localPath)
					if err == nil {
						size := fi.Size()
						currentChan <- file.UpdateCurrent(size)
						if size >= totalSize {
							return
						}
					}
				}
			}()

			// goph SCP download
			err = client.Raw().Download(remotePath, localPath)
			if err != nil {
				// fmt.Printf("failed to download %s: %w", file, err)
				// continue
			}

			finishedChan <- i + 1
		}
	}()
	// progress UI
	ui.RunProgress(len(files), currentChan, finishedChan)
	fmt.Println("All files downloaded")

	close(currentChan)
	close(finishedChan)
	return nil
}
