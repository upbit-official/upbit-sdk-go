English | [한국어](./README_KR.md)

# Example Code

This directory contains examples that demonstrate the key features of the Upbit Go SDK in scenario-based flows.

Each example is written based on real trading workflows and runs in **Dry run mode** by default.
In dry run mode, only read operations are performed — write operations such as placing orders or making withdrawals are skipped.

> **Note:** Each file is an independent `package main` program. Run with `go run examples/<file>.go`. Per-endpoint examples are organized in feature-specific subdirectories within this directory (e.g., [`examples/orders/`](./orders/) for order-related examples). For Korean scenario examples, see [`README_KR.md`](./README_KR.md).

---

## Quick Start

If you are new, we recommend going through the examples in this order:

1. [`quotation.go`](./quotation.go) — Quickly explore market data, candles, trades, and orderbook without authentication
2. [`indicators.go`](./indicators.go) — Learn how to calculate investment indicators using Quotation data
3. [`orders.go`](./orders.go) — Walk through order creation, lookup, and cancellation in Dry run first

### Examples that do not require authentication

The following examples can be run immediately without an API key:
- [`quotation.go`](./quotation.go)
- [`indicators.go`](./indicators.go)

### Examples that require authentication

The following examples require an Upbit API key:
- [`orders.go`](./orders.go)
- [`orders_validate.go`](./orders_validate.go)
- [`deposits.go`](./deposits.go)
- [`withdrawals.go`](./withdrawals.go)
- [`dca.go`](./dca.go)
- [`tp_sl.go`](./tp_sl.go)

### Prerequisites

```bash
# Issue an Upbit API key and set environment variables
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

### Before you run

- **To explore read-only behavior first**, start with `quotation.go` and `indicators.go`.
- **When reviewing examples that can place orders or make withdrawals**, verify the behavior in Dry run mode first.
- **When running with `DRY_RUN=false`**, write operations such as orders, withdrawals, and automated trading will actually execute.

---

## Example Overview

| Example | Purpose | Auth | Default Behavior |
| --- | --- | --- | --- |
| [`quotation.go`](./quotation.go) | Market data / candles / trades / orderbook | Not required | Safe read-only |
| [`indicators.go`](./indicators.go) | Calculate investment indicators such as RSI | Not required | Safe read-only |
| [`orders.go`](./orders.go) | Order creation, lookup, and cancellation flow | Required | Dry run by default; real orders possible |
| [`orders_validate.go`](./orders_validate.go) | Validate all order types via Order Creation Test API (KRW) | Required | Always safe — no real orders placed |
| [`deposits.go`](./deposits.go) | Deposit address and history management | Required | Dry run by default; some operations may have real effect |
| [`withdrawals.go`](./withdrawals.go) | Withdrawal info and withdrawal flow | Required | Dry run by default; real withdrawals require caution |
| [`dca.go`](./dca.go) | Automated recurring market-buy (DCA) | Required | Dry run by default; real orders possible |
| [`tp_sl.go`](./tp_sl.go) | Automated take-profit / stop-loss sell | Required | Dry run by default; real orders possible |

---

## Scenario Examples in Detail

### `quotation.go` — Market Data Primer

A quick introduction to the Quotation API, usable without authentication.
A great starting point if you are exploring the SDK for the first time.

**Covered features:**
- List tradable pairs (markets) and filter for caution-flagged tickers
- Fetch current price (ticker) — single, multiple, and per market
- Fetch candle (OHLCV) data — 5-minute, daily, weekly, monthly / using the `to` parameter
- Fetch recent trade history — using the `days_ago` parameter
- Fetch orderbook data

```bash
go run examples/quotation.go
```

---

### `orders.go` — Order Creation and Management

Covers the full flow from placing a limit buy order to looking it up and cancelling it.

**Covered features:**
- Check order availability (fee rate, minimum order amount)
- Fetch tick size and calculate order price
- Place a limit buy order → check status → cancel
- Check holdings and verify market sell eligibility
- List completed orders

```bash
# Dry run (default) — read-only, no orders placed
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders.go

# Live run — limit buy → lookup → cancel
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders.go
```

---

### `orders_validate.go` — Order Type Validation (KRW, TestNew)

Validates all order types against the KRW market using the Order Creation Test API.
The Order Creation Test API performs the same validation as a real order but does not create an actual order — no fees are charged and the returned UUID cannot be used for lookups or cancellations.

**Covered features:**
- Check order availability (fee rate, minimum order amount, supported order types)
- Limit buy/sell with tick-size-aware price rounding
- Market buy (price) / market sell (volume)
- Best buy/sell with IOC time-in-force
- Error handling validation (invalid market, missing price)

```bash
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders_validate.go
```

---

### `deposits.go` — Deposit Address and History Management

Covers the full range of deposit-related features, from managing deposit addresses to reviewing history.

**Covered features:**
- List deposit addresses
- Look up and generate a deposit address for a specific coin
- Check deposit availability (deposit limits, minimum deposit amount)
- Fetch deposit history — list view / single lookup by UUID or TxID
- List Travel Rule-supported exchanges (VASPs)

```bash
# Dry run (default) — read-only
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits.go

# Live run — includes deposit address generation
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits.go
```

---

### `withdrawals.go` — Withdrawal Management

Covers the full range of withdrawal-related features, from checking allowed addresses to executing a withdrawal.
Always double-check the address, network, and amount before running.

> **Warning:** Withdrawals cannot be reversed. Always verify the address and amount before executing a real withdrawal.

**Covered features:**
- List allowed withdrawal addresses
- Check deposit/withdrawal service status (BTC, ETH, USDT)
- Fetch withdrawal availability (fee, minimum amount, daily limit)
- Fetch withdrawal history — list view / single lookup
- Digital asset withdrawal and KRW fiat withdrawal

```bash
# Dry run (default) — read-only
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals.go

# Live run — fetches withdrawal availability (actual withdrawal requires code modification)
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals.go
```

---

### `dca.go` — DCA Automated Recurring Buy

Implements a Dollar Cost Averaging (DCA) strategy that repeatedly buys a fixed amount at market price.
Useful for understanding how to structure a periodic buy automation.

> This example assumes a server or scheduler environment is already set up.
> You must configure a persistent automated execution environment yourself.

**Covered features:**
- Fetch current price and verify buy configuration
- Repeat market buy for a specified number of times
- Wait for order fill and poll for status
- Summarize buy results (total spent, total quantity, average price)

```bash
# Dry run (default) — fetches current price only, no orders placed
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/dca.go

# Live run
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/dca.go
```

---

### `tp_sl.go` — Take-Profit / Stop-Loss Auto Sell

Polls the current price and executes a market sell when the target price (take-profit) or stop-loss price is reached.

> This example assumes a server or scheduler environment is already set up.

**Covered features:**
- Fetch current price and calculate TP/SL prices
  - Dry run: based on current market price
  - Live run: based on average buy price of held assets
- Price monitoring via REST polling (1-second interval)
- Execute market sell when target price is reached
- Wait for fill and confirm result

```bash
# Dry run (default) — monitoring only, no orders placed
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/tp_sl.go

# Live run — auto sell based on average buy price
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/tp_sl.go
```

---

### `indicators.go` — Investment Indicator Calculation

Calculates technical investment indicators using data from the Quotation API.
Implemented in pure Go without any external libraries.

**Covered features:**
- Find the top 5 pairs by 24-hour cumulative trading volume
- RSI (Relative Strength Index) calculation
  - Based on 200 daily candles
  - Wilder's smoothing method
  - RSI interpretation (overbought / oversold / neutral)

```bash
# No authentication required
go run examples/indicators.go
```

If you want a higher-level flow overview first, read through this README, then explore the per-endpoint examples in the feature-specific subdirectories within [`examples/`](./).
