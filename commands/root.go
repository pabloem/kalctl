package commands

import (
	"fmt"
	"os"

	"github.com/pabloem/kalctl/auth"
	"github.com/pabloem/kalctl/commands/base"
	"github.com/pabloem/kalctl/commands/impl"
	"github.com/pabloem/kalctl/output"
	"github.com/pabloem/kalctl/reqs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var RootNs = base.NewNamespace(
	"kalctl",
	"Kalctl is a CLI tool for interfacing with the Kalshi Prediction Market APIs",
	impl.EventsNamespace,
	impl.MarketsNamespace,
	impl.TradesNamespace,
	base.NewNamespace(
		"exchange",
		"Query the schedule and announcements from the exchange",
		base.NewCommand(
			"get-schedule",
			"Get the schedule of events from the exchange",
			reqs.HttpRequestTemplate{
				Method: reqs.GET,
				Path:   "trade-api/v2/exchange/schedule",
			},
		),
		base.NewCommand(
			"get-announcements",
			"Get the announcements from the exchange",
			reqs.HttpRequestTemplate{
				Method: reqs.GET,
				Path:   "trade-api/v2/exchange/status",
			},
		),
	),
	base.NewNamespace(
		"portfolio",
		"Query and manage your portfolio. This includes creating new orders and querying existing ones",
		base.NewCommand(
			"get-balance",
			"Get the balance of your account in USD",
			reqs.HttpRequestTemplate{
				Method: reqs.GET,
				Path:   "trade-api/v2/portfolio/balance",
			},
		),
		base.NewNamespace(
			"orders",
			"Query and manage your orders",
			base.NewCommand(
				"list",
				"List your open orders",
				reqs.HttpRequestTemplate{
					Method: reqs.GET,
					Path:   "trade-api/v2/portfolio/orders",
				}, // TODO: Add arguments
			),
		),
	),
	base.NewNamespace(
		"data-feed",
		"Get a data feed from the Kalshi exchange.",
	),
	base.NewNamespace(
		"auth",
		"Authenticate with the Kalshi API",
		impl.NewCustomRunCommand(
			"login",
			"Authenticate with the Kalshi API",
			func(args base.CommandArgs) error {
				_, permAuth := args.KwArgs["perm"]
				return auth.RunKalshiAuth(permAuth)
			},
			base.Argument{
				Name: "perm",
				Desc: "Store authentication credentials permanently. " +
					"With this flag, you will not need to re-authenticate on subsequent runs;" +
					"however, anyone with access to your user directory can access your credentials.",
				Position: -1,
				Required: false,
			},
		),
	),
	impl.NewCustomRunCommand("init", "Configure kalctl autocomplete",
		func(args base.CommandArgs) error {
			fmt.Println(output.GetFormatter().AttributeDescription("Add the following to your .bashrc or .zshrc file:"))
			fmt.Println(output.GetFormatter().CommandResult(impl.AUTOCOMPLETE_SCRIPT))
			return nil
		}),
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
				fmt.Println(formatter.Attribute("--" + arg.Name))
			}
			return
		}
		fmt.Println("Arguments:")
		for _, arg := range typedElm.Arguments() {
			fmt.Println(formatter.Attribute("--" + arg.Name))
			fmt.Println(formatter.AttributeDescription(arg.Desc))
		}
	}
}

func executeCommand(parsed base.CommandArgs) error {
	_, dashHelp := parsed.KwArgs["help"]
	_, shortHelp := parsed.KwArgs["short"]
	if len(parsed.Args) == 0 {
		printHelp(RootNs, shortHelp)
		return nil
	}
	var curElm base.Element = RootNs
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
			parsed.Args = parsed.Args[len(parsed.Args):]
			return typedElm.Run(parsed)
		}
	}
	return nil
}

func RunCommand(args []string) error {
	args = args[1:]
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

	return executeCommand(parsed)
}
