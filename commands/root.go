package commands

import (
	"fmt"

	"github.com/pabloem/kalctl/auth"
	"github.com/pabloem/kalctl/commands/base"
	"github.com/pabloem/kalctl/commands/impl"
	"github.com/pabloem/kalctl/output"
	"github.com/pabloem/kalctl/reqs"
)

var RootNs = base.NewNamespace(
	"kalctl",
	"Kalctl is a CLI tool for interfacing with the Kalshi Prediction Market APIs",
	impl.MarketsNamespace,
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
