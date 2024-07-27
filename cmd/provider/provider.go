/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package provider

import (
	"github.com/spf13/cobra"
	"turtle/cmd/provider/google"
)

// providerCmd represents the provider command
var ProviderCmd = &cobra.Command{
	Use:   "provider",
	Short: "Provider contains the cloud storage provider that are supported",
	Long:  "Provider contains the cloud storage provider that are supported",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func addSubCommand() {
	ProviderCmd.AddCommand(google.GoogleCmd)
}

func init() {
	addSubCommand()
}
