package main

import (
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"os"
	"slices"
	"text/tabwriter"

	"github.com/TrueHopolok/braincode-/server/config"
)

// On command encounter in the os.Stdin, the function will be executed
type Instruction struct {
	command  string
	helptext string
  function func(chan bool)
}

// Contain all instructions that can be accessed via console
var Instructions = []Instruction{
	{
		"stop",
    "alert quitChannel, thus stopping the process (should not, fix main function if that happens)",
		func(quitChan chan bool) {
			quitChan <- true
		},
	},
}

func init() {
	Instructions = append(Instructions, Instruction{
		command:  "help",
		helptext: "print description of all available commands",
		function: func(chan bool) { fmt.Println(commandHelpText()) },
	})

	slices.SortFunc(Instructions, func(l, r Instruction) int {
		return cmp.Compare(l.command, r.command)
	})
}

// Wait for the input in os.Stdin.
// Check if inputed string is one of the commands in the intructions slice.
// If it is, the function of that instruction is executed.
func ConsoleHandler(quitChan chan bool) error {
	if !config.Get().EnableConsole {
		return fmt.Errorf("Console is blocked by config parameters")
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	fmt.Println("Waiting for user input:")
	for scanner.Scan() && scanner.Err() == nil {
		// if multiarguments will be needed: use strings.Fields or flag package
		// if faster checker required, use search tree for string
		// if required auto correct use spell checker package
		request := scanner.Text()
		found := false
		for _, instruct := range Instructions {
			if instruct.command == request {
				instruct.function(quitChan)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Invalid command, try again...")
		}
	}
	return scanner.Err()
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
