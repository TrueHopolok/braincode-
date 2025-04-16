package main

import (
	"bufio"
	"fmt"
	"os"
)

type Instruction struct {
	command string
	function func()
}

var Instructions = []Instruction{
	{"stop", StopServer},
}

func ConsoleHandler() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
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