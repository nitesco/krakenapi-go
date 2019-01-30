package krakenapi_test

import (
	"fmt"
	"github.com/crankykernel/krakenapi-go"
	"log"
)

func ExampleSocket() {
	ws, err := krakenapi.OpenWebSocket()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for first message and check the status.
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

	// Subscribe to some tickers.
	ws.SubscribeTicker("XBT/USD", "ETH/USD", "XLM/USD")
	ws.SubscribeTicker("LTC/USD")

	// Listen for messages.
	for {
		// We first read in a raw message. An error here is a socket level
		// error.
		payload, err := ws.Next()
		if err != nil {
			log.Fatalf("socket error: %+v\n", err)
		}

		// Decoding is a separate step, to determine if an error was a socket
		// level error or a decoding error.
		decoded, err := ws.Decode(payload)
		if err != nil {
			log.Fatalf("decode error: %+v\n", err)
		}

		// Based on the type of message received, different types may be
		// returned from the decode method above.
		switch v := decoded.(type) {
		case krakenapi.Ticker:
			fmt.Printf("Ticker: %+v\n", v)
		case krakenapi.EventMessage:
			fmt.Printf("Event: %+v\n", v)
		}
	}
}
