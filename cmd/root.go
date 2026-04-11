package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rioctl",
	Short: "RoboRIO SSH utility",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntP("team", "t", 0, "FRC team number")
	rootCmd.PersistentFlags().StringP("user", "u", "admin", "SSH user")
	rootCmd.PersistentFlags().StringP("password", "p", "", "SSH password (default: empty)")
}
