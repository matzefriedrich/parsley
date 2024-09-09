package commands

import (
	"context"
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/spf13/cobra"
)

type versionCommand struct {
	use            abstractions.CommandName `flag:"version" short:"Show the current Parsley CLI version"`
	CheckForUpdate bool                     `flag:"check-update" usage:"Checks for available updates and prints the update command"`
}

func (v *versionCommand) Execute() {

	appVersion, appVersionErr := utils.ApplicationVersion()
	if appVersionErr == nil {
		fmt.Printf("Parsley CLI v%s\n", appVersion.String())
	}

	if v.CheckForUpdate == false {
		return
	}

	githubClient := utils.NewGitHubApiClient()
	release, err := githubClient.QueryLatestReleaseTag(context.Background())
	if err != nil {
		return
	}

	releaseVersion, releaseVersionErr := release.TryParseVersionFromTag()
	if appVersionErr == nil && releaseVersionErr == nil {
		if appVersion.LessThan(*releaseVersion) {

			fmt.Printf("\n"+
				"Your version of Parsley CLI is out of date!\n\n"+
				"The latest version is: v%s.\n"+
				"To update run the following command: "+
				"go install github.com/matzefriedrich/parsley/cmd/parsley-cli@v%s\n\n", releaseVersion.String(), releaseVersion.String())

			fmt.Printf("More information about the release %s is available at:\n%s\n", release.Name, release.HtmlUrl)

		} else if appVersion.Equal(*releaseVersion) {

			fmt.Printf("\n" +
				"You are using the latest version of Parsley CLI.\n\n")

		}
	}
}

var _ pkg.TypedCommand = (*versionCommand)(nil)

func NewVersionCommand() *cobra.Command {
	command := &versionCommand{}
	return pkg.CreateTypedCommand(command)
}
