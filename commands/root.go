package commands

import (
	"github.com/pabloem/kalctl/auth"
	"github.com/pabloem/kalctl/commands/base"
	"github.com/pabloem/kalctl/commands/impl"
	"github.com/pabloem/kalctl/reqs"
)

var RootNs = base.NewNamespace(
	"kalctl",
	"Kalctl is a CLI tool for interfacing with the Kalshi Prediction Market APIs",
	base.NewNamespace(
		"markets",
		"Interact with the markets",
		base.NewNamespace(
			"events",
			"Query the events available in the markets",
			base.NewCommand(
				"list",
				"List the events available in the markets",
				reqs.HttpRequestTemplate{
					Method: reqs.GET,
					Path:   "trade-api/v2/markets/events",
				},
				base.Argument{
					Name:     "limit",
					Desc:     "Limit the number of events returned",
					Position: -1,
					Required: false,
				},
				base.Argument{
					Name:     "series",
					Desc:     "Filter by series",
					Position: 0,
					Required: false,
				},
			),
			base.NewCommand(
				"get",
				"Get the details of a specific event",
				reqs.HttpRequestTemplate{
					Method: reqs.GET,
					Path:   "trade-api/v2/markets/events/{event}",
				},
				base.Argument{
					Name:     "event",
					Desc:     "The ticker of the event to retrieve",
					Position: 0,
					Required: true,
				},
			),
		),
	),
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
	impl.NewCustomRunCommand("init", "Configure kalctl, including autocomplete and authentication",
		func(args base.CommandArgs) error {
			return nil // TODO: Implement
		}),
)
