package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildVersion string
var buildCommit string
var buildDate string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ferrariprocessor",
	Long:  `All software has versions. This is ferrariprocessor's.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("ferrariprocessor job processor v%s-%s Build date:%s\n", buildVersion, buildCommit, buildDate)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
