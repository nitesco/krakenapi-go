# krakenapi-go

A Kraken API for Go

## REST API

REST API support is not included yet. The REST API features from
https://gitlab.com/crankykernel/cryptotrader will be ported over
soon.

## WebSocket Support

Currently supports the following websocket API features:

* Application ping
* Tickers
* OHLC
* Spread

## WebSocket Example

https://github.com/crankykernel/krakenapi-go/blob/master/example_socket_test.go