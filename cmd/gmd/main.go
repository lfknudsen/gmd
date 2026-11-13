package main

import (
	"fmt"
	"os"

	"gmd/pkg"
)

func printHelp(exitCode int) {
	fmt.Println("Usage: gmd <filename>")
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) == 0 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp(0)
	}

	buffer, err := os.ReadFile(os.Args[1])
	if err != nil {
		_, _ = fmt.Printf("Could not open the file '%s'\n", os.Args[1])
		os.Exit(1)
	}

	str := pkg.RenderMarkdown(buffer)
	fmt.Print(str)

}
