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

package main

import (
	"fmt"
	"github.com/crankykernel/krakenapi-go"
	"log"
)

func main() {
	ws, err := krakenapi.OpenSocket()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for first message.
	payload, err := ws.Next()
	if err != nil {
		log.Fatal(err)
	}
	event, err := ws.Decode(payload)
	if status, ok := event.(krakenapi.EventMessage); !ok {
		log.Printf("Received unexpected message type on connect.")
	} else {
		fmt.Printf("Event: %s; Status: %s\n", status.Event, status.Status)
	}

	ws.SubscribeTicker("XBT/USD", "ETH/USD", "XLM/USD")
	ws.SubscribeOHLC(krakenapi.Interval_1m, "XBT/USD")
	ws.SubscribeSpread("XBT/USD")

	for {
		payload, err := ws.Next()
		if err != nil {
			// An error here is a socket level error. We probably need to reconnect.
			log.Fatalf("socket error: %+v\n", err)
		}

		// Decoding is a separate step, to determine if an error was a socket
		// level error or a decoding error.
		decoded, err := ws.Decode(payload)
		if err != nil {
			log.Fatalf("decode error: %+v\n", err)
		}

		switch v := decoded.(type) {
		case krakenapi.EventMessage:
			fmt.Printf("Event: %+v\n", v)
		case krakenapi.Ticker:
			fmt.Printf("Ticker: %+v\n", v)
		case krakenapi.OHLC:
			fmt.Printf("OHLC: %+v\n", v)
		case krakenapi.Spread:
			fmt.Printf("Spread: %+v\n", v)
		default:
			fmt.Printf("Unknown type: %+v\n", v)
		}
	}
}
