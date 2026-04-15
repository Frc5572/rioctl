package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"rioctl/internal/transfer"
	"rioctl/internal/ui"
	"rioctl/internal/utils"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Log operations",
}
var defaultSource = "/media/sda1"
var logsPullCmd = &cobra.Command{
	Use:   "pull [files...]",
	Short: "Pull logs from RoboRIO",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getSSHClient(cmd)
		if err != nil {
			return err
		}
		defer client.Close()

		source, _ := cmd.Flags().GetString("source")
		// If no dir provided → prompt
		if source == "" {
			input, err := ui.RunTextInput(defaultSource)
			if err != nil {
				return err
			}

			if input == "" {
				source = defaultSource
			} else {
				source = input
			}
		}

		dest, _ := cmd.Flags().GetString("dest")
		eventCode, _ := cmd.Flags().GetString("event")
		practice, _ := cmd.Flags().GetBool("practice")
		qual, _ := cmd.Flags().GetBool("qual")
		elim, _ := cmd.Flags().GetBool("elim")

		var files []string

		// If args provided → use them directly
		if len(args) > 0 {
			files = args
		} else {
			fmt.Println("Fetching file list...")
			allFiles, err := client.ListFiles(source)
			if err != nil {
				return err
			}
			if eventCode != "" || practice || qual || elim {
				allFiles, err := client.ListFiles(source)
				if err != nil {
					return err
				}
				var totalSize int64
				for _, file := range allFiles {
					if eventCode != "" && file.Event != eventCode {
						continue
					}
					if !practice && !elim && !qual {
						files = append(files, file.Name)
						totalSize += file.Size
					} else if practice && file.MatchType == "p" {
						files = append(files, file.Name)
						totalSize += file.Size
					} else if qual && file.MatchType == "q" {
						files = append(files, file.Name)
						totalSize += file.Size
					} else if elim && file.MatchType == "e" {
						files = append(files, file.Name)
						totalSize += file.Size
					}
				}

				fmt.Printf("⚠️ This will pull %d files (%s). Continue? (yes/no): ", len(files), utils.Humanize(totalSize))

				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')

				if strings.ToLower(text) != "yes\n" {
					fmt.Println("Aborted.")
					return nil
				}

			} else {
				// Launch TUI selector
				selected, err := ui.RunFilePicker(allFiles)
				if err != nil {
					return err
				}

				files = selected
			}
		}

		fmt.Println("Downloading selected files...")
		return transfer.PullFiles(client, source, dest, files)
	},
}

func init() {
	logsPullCmd.Flags().String("source", "", "Source directory")
	logsPullCmd.Flags().String("dest", "./logs", "Destination directory")
	logsPullCmd.Flags().String("event", "", "Event Code to download Event only logs")
	logsPullCmd.Flags().Bool("practice", false, "Pull Practice Match Logs")
	logsPullCmd.Flags().Bool("qual", false, "Pull Qual Match Logs")
	logsPullCmd.Flags().Bool("elim", false, "Pull Elim Match Logs")

	logsCmd.AddCommand(logsPullCmd)
	rootCmd.AddCommand(logsCmd)
}
