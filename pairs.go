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

package krakenapi

import "strings"

var restPairMap map[string]string
var websocketPairMap map[string]string

func init() {
	initRestPairs()
	initWebSocketPairs()
}

func initRestPairs() {
	restPairMap = make(map[string]string)

	restPairMap["XBTUSD"] = "XXBTZUSD"
	restPairMap["BTCUSD"] = "XXBTZUSD"

	restPairMap["XLMUSD"] = "XXLMZUSD"
	restPairMap["XLMBTC"] = "XXLMXXBT"
	restPairMap["XLMXBT"] = "XXLMXXBT"

	restPairMap["DASHUSD"] = "XDASHZUSD"
	restPairMap["DASHXBT"] = "XDASHXXBT"

	restPairMap["XMRUSD"] = "XXMRZUSD"
	restPairMap["XMRXBT"] = "XXMRXXBT"

	restPairMap["LTCXBT"] = "XLTCXXBT"
	restPairMap["LTCBTC"] = "XLTCXXBT"
	restPairMap["LTCUSD"] = "XLTCZUSD"
}

func initWebSocketPairs() {
	websocketPairMap = make(map[string]string)

	websocketPairMap["BTCUSD"] = "XBT/USD"
	websocketPairMap["XMRUSD"] = "XMR/USD"
	websocketPairMap["DASHUSD"] = "DASH/USD"

	websocketPairMap["LTCXBT"] = "LTC/XBT"
	websocketPairMap["LTCBTC"] = "LTC/XBT"
}

// Given a common pair naming, return the Kraken format for the REST API.
func RestPair(input string) string {
	pair := strings.Replace(input, "/", "", 1)
	restPair, ok := restPairMap[pair]
	if ok {
		return restPair
	}
	return input
}

// Given a common pair naming, return the Kraken format for websockets.
func WebSocketPair(input string) string {
	pair := strings.Replace(input, "/", "", 1)
	wsPair, ok := websocketPairMap[pair]
	if ok {
		return wsPair
	}
	return input
}
