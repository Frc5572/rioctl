package cmd

import (
	"fmt"
	"rioctl/internal/sshclient"

	"github.com/spf13/cobra"
)

func teamToIP(team int) string {
	return fmt.Sprintf("10.%d.%d.2", team/100, team%100)
}

func getSSHClient(cmd *cobra.Command) (*sshclient.Client, error) {
	team, _ := cmd.Root().PersistentFlags().GetInt("team")
	user, _ := cmd.Root().PersistentFlags().GetString("user")
	pass, _ := cmd.Root().PersistentFlags().GetString("password")

	if team == 0 {
		return nil, fmt.Errorf("team number is required")
	}

	host := teamToIP(team)

	fmt.Printf("Connecting to %s as %s...\n", host, user)

	return sshclient.New(host, user, pass)
}
