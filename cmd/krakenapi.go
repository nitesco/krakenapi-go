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
	"github.com/nitesco/krakenapi-go"
	"github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatal("error: no command provided")
	}

	switch os.Args[1] {
	case "websocket":
		RunWebSocket()
	case "get":
		RunGet(os.Args[2:])
	case "post":
		RunPost(os.Args[2:])
	default:
		log.Fatalf("error: unknown command: %s", os.Args[1])
	}
}

func RunGet(args []string) {
	if len(args) < 1 {
		log.Fatal("error: not enough arguments: an endpoint is required")
	}
	client, err := krakenapi.NewRestClient("", "")
	if err != nil {
		log.Fatalf("error: failed to create rest client: %v", err)
	}
	response, err := client.Get(args[0])
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

func RunPost(args []string) {
	flags := pflag.NewFlagSet("get", pflag.ExitOnError)
	apiKey := flags.String("api-key", "", "API key")
	apiSecret := flags.String("api-secret", "", "API secret")

	flags.Parse(args)
	args = flags.Args()

	if len(args) == 0 {
		log.Fatal("error: no path specified")
	}

	path := args[0]

	params := map[string]interface{}{}
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		params[parts[0]] = parts[1]
	}

	client, err := krakenapi.NewRestClient(*apiKey, *apiSecret)
	if err != nil {
		log.Fatalf("error: failed to create rest client: %v", err)
	}
	response, err := client.Post(path, params)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}

func RunWebSocket() {
	ws, err := krakenapi.OpenWebSocket()
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
	ws.SubscribeSpread("XXBTZUSD")
	//ws.SubscribeBook("XBT/USD")

	ws.Ping(0)

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
