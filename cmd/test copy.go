package cmd

import (
	"rioctl/internal/ui"

	"github.com/spf13/cobra"
)

var test2CMD = &cobra.Command{
	Use:   "test2",
	Short: "test operations",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		allFiles := []string{"asdf", "asdf"}
		_, err := ui.RunFilePicker(allFiles)
		return err
	},
}

func init() {

	rootCmd.AddCommand(test2CMD)
}
