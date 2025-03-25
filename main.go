package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Shivam! Your Go application is running in a Docker container.")
	})

	port := "8080"
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	myEnvVar := os.Getenv("MY_ENV_VAR")
	fmt.Println(myEnvVar) // Prints: Hello, World!

}
