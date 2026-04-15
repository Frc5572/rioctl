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

	hostOptions := []string{
		teamToIP(team),
		fmt.Sprintf("roboRIO-%d-FRC.local", team),
		"172.22.11.2",
	}
	for _, host := range hostOptions {
		fmt.Printf("Connecting to %s as %s...\n", host, user)
		client, err := sshclient.New(host, user, pass)
		if err == nil {
			return client, nil
		}
		fmt.Printf("Connection to %s failed\n", host)
	}
	return nil, fmt.Errorf("Could not connect to RoboRio on any IP/Hostname")
}
