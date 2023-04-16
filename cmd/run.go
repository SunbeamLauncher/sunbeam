package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/pomdtr/sunbeam/internal"
	"github.com/pomdtr/sunbeam/types"
	"github.com/pomdtr/sunbeam/utils"
	"github.com/spf13/cobra"
)

func NewCmdRun(extensionDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "run <page>",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			extension, err := ListExtensions(extensionDir)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return extension, cobra.ShellCompDirectiveDefault
		},
		Short: "Run page from file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			extensions, err := ListExtensions(extensionDir)
			if err != nil {
				return err
			}

			for _, extension := range extensions {
				if extension != args[0] {
					continue
				}
				var cmdArgs []string
				cmdArgs = append(cmdArgs, path.Join(extensionDir, extension, extensionBinaryName))
				cmdArgs = append(cmdArgs, args[1:]...)
				return Draw(internal.NewCommandGenerator(&types.Command{
					Args: cmdArgs,
				}))
			}

			if repository, err := utils.RepositoryFromString(args[0]); err == nil {
				page := types.Page{
					Type:  types.FormPage,
					Title: "Install Extension?",
					SubmitAction: &types.Action{
						Type:  types.RunAction,
						Title: "Install",
						Command: &types.Command{
							Args: []string{os.Args[0], "extension", "install", "--open", "${input:repository}"},
						},
						Inputs: []types.Input{
							{
								Name:        "repository",
								Type:        types.TextFieldInput,
								Title:       "Repository",
								Placeholder: "Repository",
								Default:     repository.FullName(),
							},
						},
						OnSuccess: types.PushOnSuccess,
					},
				}
				return Draw(func() (*types.Page, error) {
					return &page, nil
				})
			}

			if _, err := os.Stat(args[0]); err == nil {
				return Draw(internal.NewCommandGenerator(&types.Command{
					Args: args,
				}))
			}

			if _, err := exec.LookPath(args[0]); err == nil {
				return Draw(internal.NewCommandGenerator(&types.Command{
					Args: args,
				}))
			}

			return fmt.Errorf("file or command not found: %s", args[0])
		},
	}

	cmd.Flags().String("on-success", "push", "action to trigger when the command is successful")

	return cmd
}
