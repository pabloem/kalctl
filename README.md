# Kalctl - a CLI for the Kalshi.com API

`kalctl` is a CLI to interact with the Kalshi API. After installing the CLI, you can use it to get information about events, markets, trades, and more.

**Usage**

First you must authenticate to Kalshi.

```bash
kalctl auth login --perm
```

After authenticating, you can query events, markets, etc:

```bash
kalctl events list
```

```bash
kalctl events get $EVENT_TICKER
```

```bash
kalctl markets list
```

```bash
kalctl markets get $MARKET_TICKER
```

Use `--help` to get more information about the commands.


**Output formats**

- ✅ JSON (the raw output from the Kalshi API)
- 🚧 YAML

**APIs supported**

|Status|API|Command|
|---|---|---|
|✅|[market/GetEvents](https://trading-api.readme.io/reference/getevents)|`kalctl events list`|
|✅|[market/GetEvent](https://trading-api.readme.io/reference/getevent)|`kalctl events get`|
|✅|[market/GetMarkets](https://trading-api.readme.io/reference/getmarkets)|`kalctl markets list`|
|✅|[market/GetMarket](https://trading-api.readme.io/reference/getmarket)|`kalctl markets get`|
|✅|[market/GetMarketOrderbook](https://trading-api.readme.io/reference/getmarketorderbook)|`kalctl markets orderbook get`|
|✅|[market/GetMarketCandlesticks](https://trading-api.readme.io/reference/getmarketcandlesticks)|`kalctl markets candlesticks get`|
|✅|[market/GetTrades](https://trading-api.readme.io/reference/gettrades)|`kalctl trades list`|
|✅|[market/GetSeries](https://trading-api.readme.io/reference/getseries)|`kalctl events series get`|
|✅|[exchange/GetExchangeSchedule](https://trading-api.readme.io/reference/getexchangeschedule)|`kalctl exchange schedule get`|
|✅|[exchange/GetExchangeAnnouncements](https://trading-api.readme.io/reference/getexchangeannouncements)|`kalctl exchange announcements get`|
|🚧|[portfolio/CreateOrder](https://trading-api.readme.io/reference/createorder)|`kalctl portfolio orders create`|
|🚧|[portfolio/AmendOrder](https://trading-api.readme.io/reference/amendorder)|`kalctl portfolio orders amend`|
|🚧|[portfolio/DecreaseOrder](https://trading-api.readme.io/reference/decreaseorder)|`kalctl portfolio orders decrease`|
|🚧|[portfolio/CancelOrder](https://trading-api.readme.io/reference/cancelorder)|`kalctl portfolio orders cancel`|
|🚧|[portfolio/GetOrders](https://trading-api.readme.io/reference/getorders)|`kalctl portfolio orders list`|
|🚧|[portfolio/GetOrder](https://trading-api.readme.io/reference/getorder)|`kalctl portfolio orders get`|
|🚧|[portfolio/GetFills](https://trading-api.readme.io/reference/getfills)|`kalctl portfolio fills get`|
|🚧|[portfolio/BatchCreateOrders](https://trading-api.readme.io/reference/batchcreateorders)|`kalctl portfolio orders batch create`|
|🚧|[portfolio/BatchCancelOrders](https://trading-api.readme.io/reference/batchcancelorders)|`kalctl portfolio orders batch cancel`|
|🚧|[portfolio/GetPositions](https://trading-api.readme.io/reference/getpositions)|`kalctl portfolio positions list`|
|🚧|[portfolio/GetPortfolioSettlements](https://trading-api.readme.io/reference/getportfoliosettlements)|`kalctl portfolio settlements list`|
|🚧|[portfolio/GetPortfolioRestingOrderTotalValue](https://trading-api.readme.io/reference/getportfoliorestingordertotalvalue)|`kalctl portfolio restingordertotalvalue`|


**Using Kalctl??**

If you pick up Kalctl, please let me know by starring this repo.

Pull Requests and Issue Reports are welcome!

