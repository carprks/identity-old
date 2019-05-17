package main

import (
	"fmt"
	"github.com/carprks/identity/src"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func _main(args []string) int {
	// Development
	if len(args) >= 1 {
		if args[0] == "localDev" {
			err := godotenv.Load()
			if err != nil {
				fmt.Println(fmt.Errorf("godotenv err: %v", err))
			}
			fmt.Println("Running LocalDev")
		}
	}

	port := "80"
	if len(os.Getenv("PORT")) >= 2 {
		port = os.Getenv("PORT")
	}

	// Start Server
	fmt.Println(fmt.Sprintf("Starting Server on Port :%s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), src.Routes()); err != nil {
		fmt.Println(fmt.Sprintf("HTTP: %v", err))
		return 1
	}

	fmt.Println("Died but nicely")
	return 0
}

func main() {
	os.Exit(_main(os.Args[1:]))
}
