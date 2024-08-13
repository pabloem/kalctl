package impl

import (
	"github.com/pabloem/kalctl/commands/base"
	"github.com/pabloem/kalctl/reqs"
)

var EventsNamespace = base.NewNamespace(
	"events",
	"Interact with the events available in Kalshi. "+
		"An event represents a single, real-world event, "+
		"and it can have multiple markets associated with it.",
	base.NewCommand(
		"list",
		"List the events available in the markets",
		reqs.HttpRequestTemplate{
			Method: reqs.GET,
			Path:   "trade-api/v2/events",
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
			Path:   "trade-api/v2/events/{event}",
		},
		base.Argument{
			Name:     "event",
			Desc:     "The ticker of the event to retrieve",
			Position: 0,
			Required: true,
		},
	),
	base.NewCommand(
		"get-series",
		"Get a list of events by series ticker. A series is a group of related events",
		reqs.HttpRequestTemplate{
			Method: reqs.GET,
			Path:   "trade-api/v2/series/{ticker}",
		},
		base.Argument{
			Name:     "ticker",
			Desc:     "The ticker of the series to retrieve",
			Position: 0,
			Required: true,
		},
	),
)

var MarketsNamespace = base.NewNamespace(
	"markets",
	"Interact with the markets available in Kalshi. A market represents "+
		"a specific contract that can be traded on the platform, with a single "+
		"Yes/No outcome.",
	base.NewCommand(
		"list",
		"List all markets available in the platform or filter by event",
		reqs.HttpRequestTemplate{
			Method: reqs.GET,
			Path:   "trade-api/v2/markets",
		},
		base.Argument{
			Name:     "limit",
			Desc:     "Limit the number of events returned. Default is 100",
			Position: -1,
			Required: false,
		},
		base.Argument{
			Name:     "event",
			Desc:     "Event ticker to filter by",
			Position: 0,
			Required: false,
		},
		base.Argument{
			Name:     "series",
			Desc:     "Series ticker to filter by",
			Position: -1,
			Required: false,
		},
		base.Argument{
			Name:     "tickers",
			Desc:     "A comma-separated list of tickers to filter by",
			Position: -1,
			Required: false,
		},
		base.Argument{
			Name: "closes-by",
			Desc: "Return only markets that close before this time. Can be in epoch " +
				"seconds, or in the format 2006-01-02T15:04:05Z",
			Position: -1,
			Required: false,
		},
		base.Argument{
			Name: "closes-after",
			Desc: "Return only markets that close after this time. Can be in epoch " +
				"seconds, or in the format 2006-01-02T15:04:05Z",
			Position: -1,
			Required: false,
		},
		base.Argument{
			Name:     "with-status",
			Desc:     "Return only markets with the specified status. Can be 'open', 'closed' or 'settled'",
			Position: -1,
			Required: false,
		},
	),
	base.NewCommand(
		"get",
		"Get the details of a specific market",
		reqs.HttpRequestTemplate{
			Method: reqs.GET,
			Path:   "trade-api/v2/markets/{ticker}",
		},
		base.Argument{
			Name:     "ticker",
			Desc:     "The ticker of the market to retrieve",
			Position: 0,
			Required: true,
		},
	),
	base.NewNamespace(
		"orderbook",
		"Interact with the orderbooks for specific markets. "+
			"An orderbook represents the current state of the market, "+
			"including the current best bid and ask prices.",
		base.NewCommand(
			"get",
			"Get the orderbook of a specific market",
			reqs.HttpRequestTemplate{
				Method: reqs.GET,
				Path:   "trade-api/v2/markets/{ticker}/orderbook",
			},
			base.Argument{
				Name:     "ticker",
				Desc:     "The ticker of the market to retrieve",
				Position: 0,
				Required: true,
			},
		),
	),
	base.NewNamespace(
		"candlesticks",
		"Query candlesticks for specific markets. A candlestick is a representation "+
			"of the market's price over a specific period of time.",
		base.NewCommand(
			"get",
			"Get the candlestick time-series of a specific market",
			reqs.HttpRequestTemplate{
				Method: reqs.GET,
				Path:   "trade-api/v2/series/{series}/markets/{ticker}/candlesticks",
			},
			base.Argument{
				Name:     "series",
				Desc:     "The ticker of the market to retrieve",
				Position: 0,
				Required: true,
			},
			base.Argument{
				Name:     "ticker",
				Desc:     "The ticker of the market to retrieve",
				Position: 1,
				Required: true,
			},
			base.Argument{
				Name: "since",
				Desc: "The time to start the candlesticks from. Can be in epoch " +
					"seconds, or in the format 2006-01-02T15:04:05Z",
				Position: -1,
				Required: true,
			},
			base.Argument{
				Name: "until",
				Desc: "The time to end the candlesticks. Default is the current time." +
					" Can be in epoch seconds, or in the format 2006-01-02T15:04:05Z",
				Position: -1,
				Required: false,
			},
			base.Argument{
				Name: "period",
				Desc: "The period of the candlesticks. Can be 1m, 1h or 1d for " +
					"1 minute, 1 hour or 1 day respectively",
				Position: -1,
				Required: true, // TODO: Add a period?
			},
		),
	),
)

var TradesNamespace = base.NewNamespace(
	"trades",
	"Query trades in Kalshi. Can query trades for all markets or for a"+
		"specific market",
	base.NewCommand(
		"list",
		"Get the trades for all or a subset of markets",
		reqs.HttpRequestTemplate{
			Method: reqs.GET,
			Path:   "trade-api/v2/markets/trades",
		},
		base.Argument{
			Name:     "ticker",
			Desc:     "The ticker of the market to retrieve",
			Position: 0,
			Required: false,
		},
		base.Argument{
			Name:     "limit",
			Desc:     "Limit the number of trades returned. Default is 100",
			Position: -1,
			Required: false,
		},
		base.Argument{
			Name: "since",
			Desc: "The time to start the trades from. Can be in epoch " +
				"seconds, or in the format 2006-01-02T15:04:05Z",
			Position: -1,
			Required: false,
		},
		base.Argument{
			Name: "until",
			Desc: "The time to end the trades. Default is the current time." +
				" Can be in epoch seconds, or in the format 2006-01-02T15:04:05Z",
			Position: -1,
			Required: false,
		},
	),
)
