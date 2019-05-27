// MIT License
//
// Copyright (c) 2019 Cranky Kernel
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package krakenapi provides access to the Kraken cryptocurrency exchange.
package krakenapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var WS_URL = "wss://ws.kraken.com"

type channelMeta struct {
	name string
	pair string
}

type WebSocket struct {
	Conn     *websocket.Conn
	channels map[int64]channelMeta
}

func OpenWebSocket() (*WebSocket, error) {
	conn, response, err := websocket.DefaultDialer.Dial(WS_URL, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusSwitchingProtocols {
		return nil, fmt.Errorf("failed to upgrade protocol to websocket")
	}
	return &WebSocket{
		Conn:     conn,
		channels: map[int64]channelMeta{},
	}, nil
}

func (s *WebSocket) Decode(input []byte) (interface{}, error) {
	// Return nil if there is no input data to decode.
	if input == nil || len(input) == 0 {
		return nil, nil
	}

	if input[0] == '{' {
		// Attempt to decode as a EventMessage.
		var status EventMessage
		if err := json.Unmarshal(input, &status); err != nil {
			return nil, err
		} else {
			// TODO: Handler status.ErrorMessage.
			if status.Event == "subscriptionStatus" {
				s.channels[status.ChannelID] = channelMeta{
					name: status.Subscription.Name,
					pair: status.Pair,
				}
			}
			return status, nil
		}
	} else if input[0] == '[' {
		decoded, err := decodeArray(input)
		if err != nil {
			return nil, err
		}

		channel, err := decoded[0].(json.Number).Int64()
		if err != nil {
			return nil, fmt.Errorf("failed to decode channel id: %+v", err)
		}

		meta, ok := s.channels[channel]
		if !ok {
			return nil, fmt.Errorf("failed to find type for channel: %d", channel)
		}

		switch meta.name {
		case "ticker":
			switch v := decoded[1].(type) {
			case map[string]interface{}:
				ticker, err := DecodeTicker(v)
				if err != nil {
					return nil, err
				}
				ticker.Pair = meta.pair
				return ticker, nil
			default:
				return nil, fmt.Errorf("invalid data type ticker event")
			}
		case "ohlc":
			switch v := decoded[1].(type) {
			case []interface{}:
				ohlc, err := DecodeOHLC(v)
				if err != nil {
					return nil, err
				}
				ohlc.Pair = meta.pair
				return ohlc, nil
			default:
				return nil, fmt.Errorf("invalid type for ohlc event")
			}
		case "spread":
			if v, ok := decoded[1].([]interface{}); ok {
				spread, err := DecodeSpread(v)
				if err != nil {
					return nil, err
				}
				spread.Pair = meta.pair
				return spread, nil
			} else {
				return nil, fmt.Errorf("invalid type for spread event")
			}
		default:
			return nil, fmt.Errorf("unknown channel type: %s", meta.name)
		}
	}

	return nil, nil
}

func (s *WebSocket) Next() ([]byte, error) {
	_, payload, err := s.Conn.ReadMessage()
	return payload, err
}

// Ping sends an application ping to the server. The reqId will only be
// included if non-zero.
func (s *WebSocket) Ping(reqId int) error {
	ping := map[string]interface{}{
		"event": "ping",
	}
	if reqId > 0 {
		ping["reqid"] = reqId
	}
	return s.Conn.WriteJSON(ping)
}

func (s *WebSocket) SubscribeTicker(tickers ...string) error {
	message := SubscribeMessage{
		Event: "subscribe",
		Pair:  tickers,
		Subscription: map[string]interface{}{
			"name": "ticker",
		},
	}
	return s.Conn.WriteJSON(&message)
}

func (s *WebSocket) SubscribeBook(ticker string) error {
	message := SubscribeMessage{
		Event: "subscribe",
		Pair:  []string{ticker},
		Subscription: map[string]interface{}{
			"name": "book",
		},
	}
	return s.Conn.WriteJSON(&message)
}

// Interval is a set of constants for OHLC intervals.
type Interval int

const (
	Interval_1m  Interval = 1
	Interval_5m  Interval = 5
	Interval_15m Interval = 15
	Interval_30m Interval = 30
	Interval_1h  Interval = 60
	Interval_4h  Interval = 240
	Interval_24h Interval = 1440
	Interval_7d  Interval = 10080
	Interval_15d Interval = 21600
)

func (s *WebSocket) SubscribeOHLC(interval Interval, tickers ...string) error {
	message := SubscribeMessage{
		Event: "subscribe",
		Pair:  tickers,
		Subscription: map[string]interface{}{
			"name":     "ohlc",
			"interval": interval,
		},
	}
	return s.Conn.WriteJSON(message)
}

func (s *WebSocket) SubscribeSpread(tickers ...string) error {
	message := SubscribeMessage{
		Event: "subscribe",
		Pair:  tickers,
		Subscription: map[string]interface{}{
			"name": "spread",
		},
	}
	return s.Conn.WriteJSON(message)
}

func (s *WebSocket) Close() error {
	return s.Close()
}

type EventMessage struct {
	Event        string `json:"event"`
	ChannelID    int64  `json:"channelID"`
	Status       string `json:"status"`
	Pair         string `json:"pair"`
	ErrorMessage string `json:"errorMessage"`
	RequestID    int64  `json:"reqid"`
	Subscription struct {
		Name string `json:"name"`
	}
}

type SubscribeMessage struct {
	Event        string                 `json:"event"`
	Pair         []string               `json:"pair"`
	Subscription map[string]interface{} `json:"subscription"`
}

// Ticker is the decoded representation of a ticker.
type Ticker struct {
	Pair string

	// Ask.
	Ask struct {
		Price          float64
		WholeLotVolume int64
		LotVolume      float64
	}

	// Bid.
	Bid struct {
		Price          float64
		WholeLotVolume int64
		LotVolume      float64
	}

	// Close.
	Close struct {
		Price     float64
		LotVolume float64
	}

	// Volume.
	Volume struct {
		Today       float64
		Last24Hours float64
	}

	// VWAP.
	Vwap struct {
		Today       float64
		Last24Hours float64
	}

	// Number of trades.
	Trades struct {
		Today       int64
		Last24Hours int64
	}

	// Low price.
	Low struct {
		Today       float64
		Last24Hours float64
	}

	// High price.
	High struct {
		Today       float64
		Last24Hours float64
	}

	// Open price.
	Open struct {
		Today       float64
		Last24Hours float64
	}
}

func decodeArray(input []byte) ([]interface{}, error) {
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	var decoded []interface{}
	err := decoder.Decode(&decoded)
	return decoded, err
}

// DecodeTicker decodes an array into a Ticker.
func DecodeTicker(data map[string]interface{}) (*Ticker, error) {
	var err error = nil
	var ticker *Ticker = &Ticker{}
	// Ask.
	ask, ok := data["a"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid ask")
	}
	if len(ask) < 3 {
		return nil, fmt.Errorf("not enough values in ask")
	}
	if ticker.Ask.Price, err = parseFloat(ask[0]); err != nil {
		return nil, err
	}
	if ticker.Ask.WholeLotVolume, err = ask[1].(json.Number).Int64(); err != nil {
		return nil, err
	}
	if ticker.Ask.LotVolume, err = parseFloat(ask[2]); err != nil {
		return nil, err
	}

	// Bid
	bid, ok := data["b"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid bid")
	}
	if len(bid) < 3 {
		return nil, fmt.Errorf("not enough values in bid")
	}
	if ticker.Bid.Price, err = parseFloat(bid[0]); err != nil {
		return nil, err
	}
	if ticker.Bid.WholeLotVolume, err = bid[1].(json.Number).Int64(); err != nil {
		return nil, err
	}
	if ticker.Bid.LotVolume, err = parseFloat(bid[2]); err != nil {
		return nil, err
	}

	// Close
	xclose, ok := data["c"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid close")
	}
	if ticker.Close.Price, ticker.Close.LotVolume, err = parseFloatDouble(xclose); err != nil {
		return nil, err
	}

	// Volume.
	volume, ok := data["v"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid volume")
	}
	if ticker.Volume.Today, ticker.Volume.Last24Hours, err = parseFloatDouble(volume); err != nil {
		return nil, err
	}

	// VWAP.
	vwap, ok := data["p"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid vwap")
	}
	if ticker.Vwap.Today, ticker.Vwap.Last24Hours, err = parseFloatDouble(vwap); err != nil {
		return nil, err
	}

	// Number of trades.
	trades, ok := data["t"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid trades")
	}
	if ticker.Trades.Today, err = trades[0].(json.Number).Int64(); err != nil {
		return nil, err
	}
	if ticker.Trades.Last24Hours, err = trades[1].(json.Number).Int64(); err != nil {
		return nil, err
	}

	// Low price.
	low, ok := data["l"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid low price")
	}
	if ticker.Low.Today, ticker.Low.Last24Hours, err = parseFloatDouble(low); err != nil {
		return nil, err
	}

	// High price.
	high, ok := data["h"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid high price")
	}
	if ticker.High.Today, ticker.High.Last24Hours, err = parseFloatDouble(high); err != nil {
		return nil, err
	}

	// Open price.
	open, ok := data["o"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid open price")
	}
	if ticker.Open.Today, ticker.Open.Last24Hours, err = parseFloatDouble(open); err != nil {
		return nil, err
	}

	return ticker, nil
}

// OHLC is the decoded OHLC event for a pair.
type OHLC struct {
	Pair    string
	Time    float64
	EndTime float64
	Open    float64
	High    float64
	Low     float64
	Close   float64
	VWAP    float64
	Volume  float64
	Count   int64
}

// DecodeOHLC decodes an array into an OHLC.
func DecodeOHLC(data []interface{}) (*OHLC, error) {
	var err error = nil
	var ohlc *OHLC = &OHLC{}
	if ohlc.Time, err = parseFloat(data[0]); err != nil {
		return nil, err
	}
	if ohlc.EndTime, err = parseFloat(data[1]); err != nil {
		return nil, err
	}
	if ohlc.Open, err = parseFloat(data[2]); err != nil {
		return nil, err
	}
	if ohlc.High, err = parseFloat(data[3]); err != nil {
		return nil, err
	}
	if ohlc.Low, err = parseFloat(data[4]); err != nil {
		return nil, err
	}
	if ohlc.Close, err = parseFloat(data[5]); err != nil {
		return nil, err
	}
	if ohlc.VWAP, err = parseFloat(data[6]); err != nil {
		return nil, err
	}
	if ohlc.Volume, err = parseFloat(data[7]); err != nil {
		return nil, err
	}
	if ohlc.Count, err = data[8].(json.Number).Int64(); err != nil {
		return nil, err
	}
	return ohlc, nil
}

// Spread represents a decoded spread event for a pair.
type Spread struct {
	Pair      string
	Bid       float64
	Ask       float64
	Timestamp float64
}

// DecodeSpread decodes the input instead a Spread struct.
func DecodeSpread(input []interface{}) (*Spread, error) {
	var err error = nil
	var spread *Spread = &Spread{}
	if len(input) < 3 {
		return nil, fmt.Errorf("not enough items")
	}
	if spread.Bid, err = parseFloat(input[0]); err != nil {
		return nil, err
	}
	if spread.Ask, err = parseFloat(input[1]); err != nil {
		return nil, err
	}
	if spread.Timestamp, err = parseFloat(input[2]); err != nil {
		return nil, err
	}
	return spread, nil
}

func parseFloat(input interface{}) (float64, error) {
	value, ok := input.(string)
	if !ok {
		return 0, fmt.Errorf("parseFloat: input not a string: %+v", input)
	}
	return strconv.ParseFloat(value, 64)
}

func parseFloatDouble(input []interface{}) (float64, float64, error) {
	if len(input) != 2 {
		return 0, 0, fmt.Errorf("parseFloatDouble: invalid number of elements")
	}
	a, err := parseFloat(input[0])
	if err != nil {
		return 0, 0, err
	}
	b, err := parseFloat(input[1])
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}
