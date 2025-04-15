package tag

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var devIncrementCmd = &cobra.Command{
	Use:   "dev_increment",
	Short: "Increment the latest DEV tag",
	RunE: func(cmd *cobra.Command, args []string) error {
		Dev_increment_logic()
		return nil
	},
}

func Dev_increment_logic() {

	err := FetchTags()
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	lastTag := LatestTagFetch("DEV")
	fmt.Printf("Latest and greatest tag is: %s\n", lastTag)

	newTag, _ := IncrementTag(lastTag)
	if len(newTag) == 0 {
		newTag = SetupConfigobject.CurrentWeekRelease + ".0.0-DEV.1"
	}

	fmt.Printf("::set-output name=newTag1::%s\n", newTag)
}

func init() {
	// Register the devIncrementCmd with the parent command
	TagCmd.AddCommand(devIncrementCmd)
}
