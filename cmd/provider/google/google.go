/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package google

import (
	"github.com/spf13/cobra"
	"turtle/cmd/provider/google/drive"
)

// driveCmd represents the drive command
var GoogleCmd = &cobra.Command{
	Use:   "google",
	Short: "Interact with google drive api",
	Long:  "Interact with google drive api",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func addSubCommand() {
	GoogleCmd.AddCommand(drive.DriveCmd)
}

func init() {

	addSubCommand()
}
