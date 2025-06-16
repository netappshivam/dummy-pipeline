package tag

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var promotionCmd = &cobra.Command{
	Use:   "promotional",
	Short: "Command to handle release creation logic",
	RunE: func(cmd *cobra.Command, args []string) error {
		SetupConfigobject.PromotionalFunc()
		return nil
	},
}

func (o *SetupConfig) PromotionalFunc() {
	errFetch := FetchTagsPrune()
	if errFetch != nil {
		log.Printf("Failed to fetch tags: %v", errFetch)
		return
	}

	FinalTagCheck := FetchTag(o.FinalRelease)
	FinalTagCheckFlag := false
	if FinalTagCheck == "" {
		log.Printf("Final Tag %s does not exist", o.FinalRelease)
		FinalTagCheckFlag = true
	} else {
		log.Printf("Final Tag %s exists, cannot perform operation", o.FinalRelease)
	}

	BaseTagCheck := FetchTag(o.BaseRelease)
	BaseTagCheckFlag := false
	if BaseTagCheck == "" {
		log.Printf("Base tag %s does not exist, cannot perform operation", o.BaseRelease)
	} else {
		log.Printf("Base Tag %s exists", o.BaseRelease)
		BaseTagCheckFlag = true
	}

	if o.OperationType == "FinalTag" {
		if BaseTagCheckFlag && strings.Contains(o.BaseRelease, "-RC.") && strings.Contains(o.BaseRelease, o.FinalRelease) && FinalTagCheckFlag {
			log.Println("Base release is a valid RC tag, proceeding with promotion.")
			PromotionalCreation()
		} else {
			log.Fatalf("Base release %s is not a valid RC tag or Final Release has a naming error %s", o.BaseRelease, o.FinalRelease)
			return
		}
	} else if o.OperationType == "HFFirstRelease" {
		if BaseTagCheckFlag && CheckForHFfinalName() && FinalTagCheckFlag {
			log.Println("Base release is a valid HF tag, proceeding with promotion.")
			PromotionalCreation()
		} else {
			log.Fatalf("Base release %s is not a valid HF tag or Final Release has a naming error %s", o.BaseRelease, o.FinalRelease)
			return
		}
	} else {
		log.Fatalf("Invalid operation type: %s. Operation can only be FinalTag or HFFirstRelease. ", o.OperationType)
		return
	}
}

func PromotionalCreation() {
	errFetch := FetchTagsPrune()
	if errFetch != nil {
		log.Printf("Failed to fetch tags: %v", errFetch)
		return
	}

	err := CleanWorkingDirectory()
	if err != nil {
		log.Fatalf("Error cleaning working directory: %v", err)
	}

	errGitTag := CreateGitTag(SetupConfigobject.FinalRelease, SetupConfigobject.BaseRelease)
	if errGitTag != nil {
		log.Fatalf("Error creating git tag: %v", errGitTag)
	}

	errGitPush := GitPush(SetupConfigobject.FinalRelease)
	if errGitPush != nil {
		log.Fatalf("Error pushing git tag: %v", errGitPush)
	}

	if errWrite := os.WriteFile(os.Getenv("GITHUB_OUTPUT"), []byte(fmt.Sprintf("FINAL_TAG=%s\n", SetupConfigobject.FinalRelease)), 0644); errWrite != nil {
		log.Fatalf("Error writing to stdout: %v", errWrite)
	}
}

func init() {
	TagCmd.AddCommand(promotionCmd)
}
