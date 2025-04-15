package tag

import (
	"github.com/spf13/cobra"
)

var TagCmd = &cobra.Command{
	Use:   "tag",
	Short: "A conalities",
}

func init() {

	TagCmd.AddCommand(releaseCmd)
	TagCmd.AddCommand(devIncrementCmd)

}
