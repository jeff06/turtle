/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"context"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"strings"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  "",
	Run:   executeDrive,
}

func executeDrive(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithCredentialsFile("/home/jeffrey/.config/turtle-cli/turtle-cli-drive-key.json"))
	var nextPageToken string
	qfilter, _ := cmd.Flags().GetString("qfilter")
	order, _ := cmd.Flags().GetString("order")
	pageSize, _ := cmd.Flags().GetInt("pagesize")
	maxResult, _ := cmd.Flags().GetInt("maxresult")
	currentNumberOfResult := 0
	for {
		filesListCall := driveService.Files.List()

		if pageSize > 0 {
			filesListCall.PageSize(int64(pageSize))
		} else {
			fmt.Println("pagesize must be greater than zero")
			return
		}

		if order != "" {
			filesListCall.OrderBy(order)
		}

		if qfilter != "" {
			filesListCall.Q(qfilter)
		}

		if nextPageToken != "" {
			filesListCall.PageToken(nextPageToken)
		}

		fileList, err := filesListCall.Do()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, file := range fileList.Files {
			if currentNumberOfResult >= maxResult {
				return
			}
			currentNumberOfResult++
			listReturnedFile(file)
		}

		nextPageToken = fileList.NextPageToken
		if nextPageToken == "" {
			break
		}
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if char == '\x00' {
			continue
		}
	}

}

func listReturnedFile(file *drive.File) {
	var fileEntityName string = ""
	if strings.HasPrefix(file.MimeType, "application/vnd.google-apps.folder") {
		fileEntityName += "folder"
	} else {
		fileEntityName = file.MimeType
	}
	fmt.Printf("%s - %s - %s", fileEntityName, file.Id, file.Name)
	fmt.Printf("\n")
}

func init() {
	ListCmd.Flags().StringP("order", "o", "", "Order. For syntax help, visit https://pkg.go.dev/google.golang.org/api/drive/v3#FilesListCall.OrderBy")
	ListCmd.Flags().StringP("qfilter", "q", "", "Q Filter. For syntax help, visit https://developers.google.com/drive/api/guides/search-files")
	ListCmd.Flags().IntP("pagesize", "p", 10, "Number of items shown per page")
	ListCmd.Flags().IntP("maxresult", "m", 30, "Number of total items retrieved")
}
