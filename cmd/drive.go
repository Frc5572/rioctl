package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rioctl/internal/googledrive"

	"github.com/spf13/cobra"
)

var (
	driveDirectory string
	driveReauth    bool
)

var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Google Drive operations",
}

var driveUploadCmd = &cobra.Command{
	Use:   "upload [files...]",
	Short: "Upload robot logs to Google Drive",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var files []string

		if driveDirectory != "" {
			err := filepath.WalkDir(driveDirectory, func(path string, entry os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !entry.IsDir() && strings.EqualFold(filepath.Ext(path), ".wpilog") {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		files = append(files, args...)
		if len(files) == 0 {
			return fmt.Errorf("provide file paths or use --dir")
		}

		client, err := googledrive.Authorize(
			ctx, driveReauth,
		)
		if err != nil {
			return err
		}

		uploader, err := googledrive.NewUploader(ctx, client)
		if err != nil {
			return err
		}

		return uploader.UploadMany(ctx, files)
	},
}

func init() {
	driveUploadCmd.Flags().StringVar(
		&driveDirectory, "dir", "",
		"Recursively scan a directory for .wpilog files",
	)
	driveUploadCmd.Flags().BoolVar(
		&driveReauth, "reauth", false,
		"Force Google OAuth authorization",
	)

	driveCmd.AddCommand(driveUploadCmd)
	rootCmd.AddCommand(driveCmd)
}

var _ = context.Background
