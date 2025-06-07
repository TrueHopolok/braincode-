package main

import (
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"os"
	"slices"
	"text/tabwriter"

	"github.com/TrueHopolok/braincode-/server/logger"
)

// On command encounter in the os.Stdin, the function will be executed
type Instruction struct {
	command  string
	helptext string
	function func()
}

// Contain all instructions that can be accessed via console
var Instructions = []Instruction{
	{
		"stop",
		"kill the process (may cause data loss)",
		func() {
			logger.Log.Info("Console: server stopped")
			os.Exit(0)
		},
	},
}

func init() {
	Instructions = append(Instructions, Instruction{
		command:  "help",
		helptext: "print description of all available commands",
		function: func() { fmt.Println(commandHelpText()) },
	})

	slices.SortFunc(Instructions, func(l, r Instruction) int {
		return cmp.Compare(l.command, r.command)
	})
}

// Wait for the input in os.Stdin.
// Check if inputed string is one of the commands in the intructions slice.
// If it is, the function of that instruction is executed.
func ConsoleHandler() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	logger.Log.Info("Console: initialized")
	fmt.Println("Waiting for user input:")
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

func commandHelpText() string {
	b := new(bytes.Buffer)
	w := tabwriter.NewWriter(b, 0, 4, 1, ' ', 0)

	fmt.Fprintln(w, "Available commands:")
	for _, cmd := range Instructions {
		fmt.Fprintf(w, "    %s\t- %s\n", cmd.command, cmd.helptext)
	}

	_ = w.Flush() // error ignored: write to bytes.Buffer never fails

	return b.String()
}
