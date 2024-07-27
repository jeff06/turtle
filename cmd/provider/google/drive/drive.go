/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package drive

import (
	"github.com/spf13/cobra"
	"turtle/cmd/provider/google/drive/list"
)

// driveCmd represents the drive command
var DriveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Interact with google drive api",
	Long:  "Interact with google drive api",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func addSubCommand() {
	DriveCmd.AddCommand(list.ListCmd)
}

func init() {
	addSubCommand()
}
