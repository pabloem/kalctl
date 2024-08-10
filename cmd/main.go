package main

import (
	"fmt"
	"os"

	"github.com/pabloem/kalctl/commands"
	"github.com/pabloem/kalctl/commands/base"
	"github.com/pabloem/kalctl/output"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func printHelp(elm base.Element, short bool) {
	formatter := output.GetFormatter()
	if !short {
		fmt.Println(formatter.Description(elm.Description()))
	}
	switch typedElm := elm.(type) {
	case base.Namespace:
		if short {
			for _, child := range typedElm.Children() {
				fmt.Println(formatter.Attribute(child.Name()))
			}
			return
		}
		// fmt.Println("Commands:")
		for _, child := range typedElm.Children() {
			fmt.Println(formatter.Attribute(child.Name()))
			fmt.Println(formatter.AttributeDescription(child.Description()))
		}
	case base.Command:
		if short {
			for _, arg := range typedElm.Arguments() {
				fmt.Println(formatter.Attribute(arg.Name))
			}
			return
		}
		fmt.Println("Arguments:")
		for _, arg := range typedElm.Arguments() {
			fmt.Println(formatter.Attribute(arg.Name))
			fmt.Println(formatter.AttributeDescription(arg.Desc))
		}
	}
}

func executeCommand(parsed base.CommandArgs) error {
	_, dashHelp := parsed.KwArgs["help"]
	_, shortHelp := parsed.KwArgs["short"]
	if len(parsed.Args) == 0 {
		printHelp(commands.RootNs, shortHelp)
		return nil
	}
	var curElm base.Element = commands.RootNs
	for i, elm := range parsed.Args {
		switch typedElm := curElm.(type) {
		case base.Namespace:
			found := false
			for _, child := range typedElm.Children() {
				if child.Name() == elm {
					curElm = child
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("could not find command %s in %s", elm, parsed.Args[:i])
			}
		case base.Command:
			if dashHelp {
				// TODO: Print help for a command better
				printHelp(typedElm, shortHelp)
				return nil
			}
			// Truncate the args to the current command
			parsed.Args = parsed.Args[i:]
			return typedElm.Run(parsed)
		}
	}
	if dashHelp {
		printHelp(curElm, shortHelp)
	} else {
		switch typedElm := curElm.(type) {
		case base.Namespace:
			printHelp(typedElm, shortHelp)
		case base.Command:
			return typedElm.Run(parsed)
		}
	}
	return nil
}

func main() {
	args := os.Args[1:]
	var parsed = base.ParseArgs(args)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if level, ok := parsed.KwArgs["log"]; ok {
		switch level {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		default:
			fmt.Println(fmt.Errorf("unknown log level %s", level))
			os.Exit(1)
		}
	}

	log.Debug().Msgf("Parsed args: %v", parsed)

	err := executeCommand(parsed)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
