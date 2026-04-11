package cmd

import (
	"os"
	"rioctl/internal/ui"

	"github.com/spf13/cobra"

	"charm.land/lipgloss/v2"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

var testCMD = &cobra.Command{
	Use:   "test",
	Short: "test operations",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		progressChan := make(chan int)
		files := []string{"asdf", "qwer", "1234"}
		// for i, _ := range files {
		// 	// fmt.Printf("Downloading File %s", file)
		// 	time.Sleep(4 * time.Second) // Pauses for 1 second
		// 	progressChan <- i + 1
		// }
		ui.RunProgress1(len(files), files)

		close(progressChan)
		os.Exit(0)
		return nil
	},
}

func init() {

	rootCmd.AddCommand(testCMD)
}
