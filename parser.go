package main

import (
	"fmt"
	"strings"
)

type Parser struct {
	// The number of expected arguments we should see while parsing
	NumberOfExpectedArguments int
	// The character length of the next argument to parse
	LengthOfNextArgument int
	// The amount of arguments parsed
	NumberOfArgumentsParsed int
}

// Parse converts an incoming byte string into usable command with its args
func Parse(command []byte) string {
	var cmd string
	args := []string{}

	parser := Parser{
		NumberOfExpectedArguments: 0,
		LengthOfNextArgument:      0,
		NumberOfArgumentsParsed:   0,
	}

	position := 0

	for position < len(command) {
		/*
			Example command:
			*3
			$3
			SET
			$3
			foo
			$3
			bar
		*/
		switch command[position] {
		// * represents the number of arguments that follow
		case '*':
			result, err := getIntArg(position+1, command)
			if err != nil {
				return err.Error()
			}

			parser.NumberOfExpectedArguments = result.Result
			args = make([]string, parser.NumberOfExpectedArguments-1)

			position += result.PositionsParsed
		// $ represents the length of the next argument
		case '$':
			result, err := getIntArg(position+1, command)
			if err != nil {
				return err.Error()
			}
			// Enforce the length of the next argument
			parser.LengthOfNextArgument = result.Result

			position += result.PositionsParsed
		case '\r', '\n':
			position += 1
		default:
			// Make sure we haven't reached here improperly through invalid argument syntax
			if parser.NumberOfExpectedArguments == 0 || parser.LengthOfNextArgument == 0 || parser.NumberOfArgumentsParsed >= parser.NumberOfExpectedArguments {
				return stringMsg("Invalid syntax")
			}

			parsedItem := string(command[position : parser.LengthOfNextArgument+position])

			// Ex: SET foo bar
			/*
					cmd = "SET"

					args = []string{
					    "foo",
					    "bar",
					}

				   When SET numberOfArgumentsParsed = 0 so the first arg is the command
				   When foo numberOfArgumentsParsed = 1 so the args[1-1]= args[0] = foo
				   When bar numberOfArgumentsParsed = 2 so the args[2-1]= args[1] = bar

			*/
			if parser.NumberOfArgumentsParsed > 0 {

				args[parser.NumberOfArgumentsParsed-1] = parsedItem
			} else {
				cmd = parsedItem
			}

			// Suppose we read SET then its length is 3 so position += 3
			position += parser.LengthOfNextArgument
			// Reset
			parser.LengthOfNextArgument = 0

			parser.NumberOfArgumentsParsed += 1
		}
	}
	return ParseCommand(cmd, args)
}

func ParseCommand(command string, args []string) string {
	cmd := strings.ToUpper(command)
	fmt.Printf("Received '%s' command\n", cmd)

	switch cmd {
	case "PING":
		return PerformPong(args)
	case "SET":
		return PerformSet(args)
	case "GET":
		return PerformGet(args)
	case "DEL":
		return PerformDel(args)
	default:
		return errorMsg("Unknown command")
	}
}
