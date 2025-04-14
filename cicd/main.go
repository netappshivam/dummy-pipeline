package main

import (
	"dummy-pipeline/cicd/dev_logic"
)

func main() {
	//if len(os.Args) < 2 {
	//	fmt.Println("Please provide an argument: 'release' or 'dev'")
	//	os.Exit(1)
	//}
	//
	//arg := os.Args[1]

	dev_logic.Dev_increment_logic()

	//switch arg {
	//case "release":
	//	release_branch_logic.Release_creation()
	//case "dev":
	//	// Call another function for the "dev" condition
	//	dev_logic.Dev_increment_logic()
	//default:
	//	fmt.Println("Invalid argument. Use 'release' or 'dev'.")
	//	os.Exit(1)
	//}
}

//comment added
