package api

import (
	"github.com/skema-dev/skemabuild/internal/pkg/http"
	"github.com/skema-dev/skemabuild/internal/pkg/pattern"
	"os"
	"path/filepath"
	"strings"

	"github.com/skema-dev/skema-go/logging"
	"github.com/skema-dev/skemabuild/internal/api"
	"github.com/skema-dev/skemabuild/internal/pkg/console"
	"github.com/skema-dev/skemabuild/internal/pkg/io"
	"github.com/spf13/cobra"
)

const (
	createDescription     = "Create API Stubs"
	createLongDescription = "skbuild api create --go_option github.com/com/test --input ./Hello1.proto -o ./stub-test"
)

func newCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: createDescription,
		Long:  createLongDescription,
		Run: func(c *cobra.Command, args []string) {
			input, _ := c.Flags().GetString("input")
			output, _ := c.Flags().GetString("output")
			stubTypes, _ := c.Flags().GetString("type")
			goOption, _ := c.Flags().GetString("go_option")

			if debug, _ := c.Flags().GetBool("debug"); debug {
				logging.Init("debug", "console")
			}

			content := getContentFromInputPath(input)
			stubs, err := generateStubsFromProto(content, stubTypes, goOption)
			console.FatalIfError(err)

			for path, stub := range stubs {
				stubFilepath := filepath.Join(output, path)
				if err = io.SaveToFile(stubFilepath, []byte(stub)); err != nil {
					panic(err)
				}
				console.Infof("%s\n", stubFilepath)
			}
		},
	}

	cmd.Flags().StringP("input", "i", "", "path of input protobuf file")
	cmd.Flags().StringP("output", "o", "./", "output path for generated stubs")
	cmd.Flags().String("go_option", "", "go_package option to be used in stub")
	cmd.Flags().StringP("type", "t", "grpc-go,openapi", "stub types to generate.")
	cmd.Flags().Bool("debug", false, "enable debug output")

	cmd.MarkFlagRequired("input")

	return cmd
}

func generateStubsFromProto(
	content string,
	stubTypes string,
	goOption string,
) (stubs map[string]string, err error) {
	stubs = make(map[string]string)
	stubTypeArr := strings.Split(stubTypes, ",")

	for _, stubType := range stubTypeArr {
		stubType = strings.TrimRight(strings.TrimLeft(stubType, " "), " ")
		var creator api.StubCreator

		switch stubType {
		case "grpc-go":
			creator = api.NewGoStubCreator(goOption)
		case "openapi":
			creator = api.NewOpenapiStubCreator()
		default:
			console.Errorf("unsupported stub type: %s", stubType)
			continue
		}
		contents, err := creator.Generate(content)
		console.FatalIfError(err)

		for k, v := range contents {
			stubFilePath := filepath.Join(stubType, k)
			stubs[stubFilePath] = v
		}
	}
	logging.Debugf("%v", stubs)
	return stubs, nil
}

func getContentFromInputPath(input string) string {
	var content string
	if pattern.IsHttpUrl(input) {
		content = http.GetTextContent(input)
	} else if pattern.IsGithubUrl(input) {
		console.Fatalf("please use the raw content link for github proto file")
	} else {
		data, err := os.ReadFile(input)
		console.FatalIfError(err, "failed to read from "+input)
		content = string(data)
	}
	return content
}
