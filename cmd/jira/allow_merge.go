package jira

import (
	"github.com/spf13/cobra"
	"log"
	"main/cmd/github"
	"os"
)

var allowMergeCmd = &cobra.Command{
	Use:   "allow_merge",
	Short: "Command to check if a PR is allowed to be merged with basic jira validations",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := allowMerge()
		if err != nil {
			return err
		}
		return nil
	},
}

type Jira struct{}

type JiraResponse struct {
	Fields struct {
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Assignee struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"assignee"`
	} `json:"fields"`
}

func allowMerge() error {

	_, credentials := GetJiraUrlCredentials()

	jiraID, err := ExtractJiraID(github.PrTitle)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}

	issue, err := GetJiraIssue(jiraID, credentials, defaultUrl)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}

	issueType := issue.Fields.Type.Name
	status := issue.Fields.Status.Name
	emailAddress := issue.Fields.Assignee.EmailAddress

	if issueType != "Story" && issueType != "Bug" {
		log.Println("Error: Issue type can be only 'Story' or 'Bug'.")
		os.Exit(1)
	}

	if status != "In Development" {
		log.Println("Error:Issue Status is not 'IN Development'.")
		os.Exit(1)
	}

	user, err := github.GetGithubUser(github.GhToken, github.PrUser)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}

	if user.Email == nil {
		log.Println("Error: Email not available for user:", github.PrUser)
		os.Exit(1)
	}

	ghEmail := *user.Email
	log.Println("GitHub PR Raised Email:", ghEmail)

	if emailAddress != ghEmail {
		log.Println("Email mismatch. Authors do not match.")
		os.Exit(1)
	} else {
		log.Println("Allowed to merge.")
	}
	return nil
}

func init() {
	JiraCmd.AddCommand(allowMergeCmd)
}
