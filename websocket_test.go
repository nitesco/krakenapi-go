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

import (
	"encoding/json"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestDecoderDicker(t *testing.T) {
	rawTicker := `[2,{"a":["3571.10000",23,"23.14437961"],"b":["3571.00000",6,"6.04250191"],"c":["3571.00000","0.01500000"],"v":["302.04621455","3263.36256626"],"p":["3571.17077","3561.39554"],"t":[655,4730],"l":["3565.80000","3545.50000"],"h":["3577.40000","3571.60000"],"o":["3571.20000","3542.30000"]}]`

	v, err := decodeArray([]byte(rawTicker))
	assert.Nil(t, err)
	assert.NotNil(t, v)

	channelId, err := v[0].(json.Number).Int64()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), channelId)
	ticker, err := DecodeTicker(v[1].(map[string]interface{}))
	assert.Nil(t, err)
	assert.Equal(t, ticker.Ask.Price, 3571.1)

	// "t":[655,4730]
	assert.Equal(t, ticker.Trades.Today, int64(655))
	assert.Equal(t, ticker.Trades.Last24Hours, int64(4730))

	// "l":["3565.80000","3545.50000"]
	assert.Equal(t, ticker.Low.Today, 3565.8)
	assert.Equal(t, ticker.Low.Last24Hours, 3545.5)
}

func TestDecodeOHLC(t *testing.T) {
	// Taken from the sandbox stream, but duplicate values altered for testing.
	input := `[4,["1548482894.938321","1548482940.000000","3569.90000","3569.80000","3569.70000","3569.60000","3569.50000","0.25000000",1]]`

	array, err := decodeArray([]byte(input))
	assert.Nil(t, err)
	assert.NotNil(t, array)

	channelId, err := array[0].(json.Number).Int64()
	assert.Nil(t, err)
	assert.Equal(t, int64(4), channelId)

	ohlc, err := DecodeOHLC(array[1].([]interface{}))
	assert.Nil(t, err)
	assert.Equal(t, ohlc.Time, 1548482894.938321)
	assert.Equal(t, ohlc.EndTime, 1548482940.0)
	assert.Equal(t, ohlc.Open, 3569.90000)
	assert.Equal(t, ohlc.High, 3569.80000)
	assert.Equal(t, ohlc.Low, 3569.70000)
	assert.Equal(t, ohlc.Close, 3569.60000)
	assert.Equal(t, ohlc.VWAP, 3569.50000)
	assert.Equal(t, ohlc.Volume, 0.25000000)
	assert.Equal(t, ohlc.Count, int64(1))
}

func TestDecodeSpread(t *testing.T) {
	input := `[3,["3554.80000","3554.90000","1548521950.519272"]]`

	array, err := decodeArray([]byte(input))
	assert.Nil(t, err)
	assert.NotNil(t, array)

	channelId, err := array[0].(json.Number).Int64()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), channelId)

	spread, err := DecodeSpread(array[1].([]interface{}))
	assert.Nil(t, err)
	assert.Equal(t, spread.Bid, 3554.8)
	assert.Equal(t, spread.Ask, 3554.9)
	assert.Equal(t, spread.Timestamp, 1548521950.519272)
}
