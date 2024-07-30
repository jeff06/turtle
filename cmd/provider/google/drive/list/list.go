/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"context"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	qFilter       string
	order         string
	pageSize      int
	maxResult     int
	preventEnter  bool
	fields        string
	error         []error
	nextPageToken string
}

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
	currentNumberOfResult := 0
	currentOption, isValidated := validateOptions(cmd)
	if !isValidated {
		return
	}
	t := table.NewWriter()
	splitField := strings.Split(currentOption.fields, ",")
	for i := range splitField {
		splitField[i] = strings.TrimSpace(splitField[i])
	}
	displayTableHeader(t, splitField)
	for {
		filesListCall := buildQuery(driveService, currentOption)

		fileList, err := filesListCall.Do()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, file := range fileList.Files {
			if currentNumberOfResult >= currentOption.maxResult {
				return
			}
			currentNumberOfResult++
			listReturnedFile(file, t, currentNumberOfResult, splitField)
		}

		currentOption.nextPageToken = fileList.NextPageToken
		if currentOption.nextPageToken == "" {
			break
		}

		if !currentOption.preventEnter {
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}
			if char == '\x00' {
				continue
			}
		}
	}
	t.Render()
}

func buildQuery(service *drive.Service, currentOption Options) *drive.FilesListCall {
	filesListCall := service.Files.List()
	if currentOption.pageSize > 0 {
		filesListCall.PageSize(int64(currentOption.pageSize))
	} else {
		fmt.Println("pagesize must be greater than zero")
		return nil
	}

	if currentOption.fields != "" {
		var currentFields []googleapi.Field = []googleapi.Field{"nextPageToken", googleapi.Field("files(" + currentOption.fields + ")")}
		filesListCall.Fields(currentFields...)
	}

	if currentOption.order != "" {
		filesListCall.OrderBy(currentOption.order)
	}

	if currentOption.qFilter != "" {
		filesListCall.Q(currentOption.qFilter)
	}

	if currentOption.nextPageToken != "" {
		filesListCall.PageToken(currentOption.nextPageToken)
	}

	return filesListCall
}

func displayTableHeader(t table.Writer, fields []string) {
	t.SetOutputMirror(os.Stdout)
	var currentRow = table.Row{"#"}
	for _, field := range fields {
		currentRow = append(currentRow, field)
	}
	t.AppendHeader(currentRow)
}

func listReturnedFile(file *drive.File, t table.Writer, l int, fields []string) {
	elementTodDisplay := ""
	var currentRow = table.Row{l}
	for _, field := range fields {
		switch field {
		case "mimeType":
			if strings.HasPrefix(file.MimeType, "application/vnd.google-apps.folder") {
				elementTodDisplay += "folder"
			} else {
				elementTodDisplay = file.MimeType
			}
			break
		case "name":
			elementTodDisplay = file.Name
			break
		case "id":
			elementTodDisplay = file.Id
			break
		case "trashed":
			elementTodDisplay = strconv.FormatBool(file.Trashed)
			break
		default:
			elementTodDisplay = "Not supported"
		}
		currentRow = append(currentRow, elementTodDisplay)
	}
	t.AppendRow(currentRow)
	fmt.Printf("\n")
}

func validateOptions(cmd *cobra.Command) (Options, bool) {
	qfilter, qfilterError := cmd.Flags().GetString("qfilter")
	order, orderError := cmd.Flags().GetString("order")
	pageSize, pageSizeError := cmd.Flags().GetInt("pagesize")
	maxResult, maxResultError := cmd.Flags().GetInt("maxresult")
	preventEnter, preventEnterError := cmd.Flags().GetBool("prevententer")
	fields, fieldsError := cmd.Flags().GetString("fields")
	currentOption := Options{
		qFilter:      qfilter,
		order:        order,
		pageSize:     pageSize,
		maxResult:    maxResult,
		preventEnter: preventEnter,
		fields:       fields,
		error:        []error{qfilterError, orderError, pageSizeError, maxResultError, preventEnterError, fieldsError},
	}

	return currentOption, true
}

func init() {
	ListCmd.Flags().StringP("order", "o", "", "Order. For syntax help, visit https://pkg.go.dev/google.golang.org/api/drive/v3#FilesListCall.OrderBy")
	ListCmd.Flags().StringP("qfilter", "q", "", "Q Filter. For syntax help, visit https://developers.google.com/drive/api/guides/search-files")
	ListCmd.Flags().IntP("pagesize", "p", 10, "Number of items shown per page")
	ListCmd.Flags().IntP("maxresult", "m", 30, "Number of total items retrieved")
	ListCmd.Flags().BoolP("prevententer", "e", false, "Prevent enter at the end of each pagesize")
	ListCmd.Flags().StringP("fields", "f", "", "Fields to be displayed. Not all field ar implemented. If one you require is not there, please open an issue. For more info, visit https://developers.google.com/drive/api/reference/rest/v3/files")
}
