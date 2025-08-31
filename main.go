/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"
	"rout/cmd"
)

func main() {
	initialCwd, err := os.Getwd()
	if err != nil {
		// Handle error, perhaps log it or exit
		os.Exit(1)
	}
	cmd.Execute(initialCwd)
}
