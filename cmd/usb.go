package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var usbCmd = &cobra.Command{
	Use:   "usb",
	Short: "USB operations",
}

var usbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List USB devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getSSHClient(cmd)
		if err != nil {
			return err
		}
		defer client.Close()

		out, err := client.Run("lsblk")
		if err != nil {
			return err
		}

		fmt.Println(out)
		return nil
	},
}

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format USB drive (FAT32)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getSSHClient(cmd)
		if err != nil {
			return err
		}
		defer client.Close()

		device, _ := cmd.Flags().GetString("device")
		if device == "" {
			return fmt.Errorf("device is required (e.g. /dev/sda1)")
		}

		fmt.Printf("⚠️ This will erase %s. Continue? (yes/no): ", device)

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		if strings.ToLower(text) != "yes\n" {
			fmt.Println("Aborted.")
			return nil
		}

		out, err := client.Run("RESTART=no /usr/local/frc/bin/frcKillRobot.sh")
		if err != nil {
			return err
		}
		fmt.Println(out)

		unmountCmd := fmt.Sprintf("umount %s", device)
		out, err = client.Run(unmountCmd)
		if err != nil {
			return err
		}
		fmt.Println(out)

		formatCmdStr := fmt.Sprintf("mkfs.vfat -F 32 %s", device)
		out, err = client.Run(formatCmdStr)
		if err != nil {
			return err
		}
		fmt.Println(out)

		mountCmdStr := fmt.Sprintf("mount %s /media/sda1", device)
		out, err = client.Run(mountCmdStr)
		if err != nil {
			return err
		}
		fmt.Println(out)

		out, err = client.Run("RESTART=yes /usr/local/frc/bin/frcKillRobot.sh")
		if err != nil {
			return err
		}
		fmt.Println(out)

		fmt.Println(out)
		return nil
	},
}

func init() {
	formatCmd.Flags().String("device", "", "Device to format")

	usbCmd.AddCommand(usbListCmd)
	usbCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(usbCmd)
}
