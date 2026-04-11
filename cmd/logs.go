package cmd

import (
	"fmt"

	"rioctl/internal/transfer"
	"rioctl/internal/ui"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Log operations",
}

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
			input, err := ui.RunTextInput("/media/sda1")
			if err != nil {
				return err
			}

			if input == "" {
				source = "/media/sda1"
			} else {
				source = input
			}
		}

		dest, _ := cmd.Flags().GetString("dest")

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

			// Launch TUI selector
			selected, err := ui.RunFilePicker(allFiles)
			if err != nil {
				return err
			}

			files = selected
		}

		fmt.Println("Downloading selected files...")
		return transfer.PullFiles(client, source, dest, files)
	},
}

func init() {
	logsPullCmd.Flags().String("source", "/media/sda1", "Source directory")
	logsPullCmd.Flags().String("dest", "./logs", "Destination directory")

	logsCmd.AddCommand(logsPullCmd)
	rootCmd.AddCommand(logsCmd)
}
