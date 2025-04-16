package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/TrueHopolok/braincode-/back-end/logger"
)

// On command encounter in the os.Stdin, the function will be executed
type Instruction struct {
	command  string
	function func()
}

// Contain all instructions that can be accessed via console
var Instructions = []Instruction{
	{
		"stop",
		func() {
			logger.Log.Info("Console: server stopped")
			logger.Stop()
			os.Exit(0)
		},
	},
}

// Wait for the input in os.Stdin.
// Check if inputed string is one of the commands in the intructions slice.
// If it is, the function of that instruction is executed.
func ConsoleHandler() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	logger.Log.Info("Console: initialized")
	for scanner.Scan() {
		// if multiarguments will be needed: use strings.Fields or flag package
		// if faster checker required, use search tree for string
		// if required auto correct use spell checker package
		request := scanner.Text()
		found := false
		for _, instruct := range Instructions {
			if instruct.command == request {
				instruct.function()
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Invalid command, try again...")
		}
	}
}
