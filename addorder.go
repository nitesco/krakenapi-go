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
	"fmt"
	"io/ioutil"
)

type OrderSide string

const OrderSideBuy OrderSide = "buy"
const OrderSideSell OrderSide = "sell"

type AddOrderOrderType string

const OrderTypeLimit AddOrderOrderType = "limit"
const OrderTypeMarket AddOrderOrderType = "market"

type AddOrderRequest struct {
	Pair         string
	Side         OrderSide
	Type         AddOrderOrderType
	Price        float64
	Volume       float64
	UserRef      int32
	ValidateOnly bool
}

type AddOrderResponse struct {
	Error  []interface{} `json:"error"`
	Result struct {
		Descr struct {
			Order string `json:"order"`
		}
		Txid []string `json:"txid"`
	} `json:"result"`
}

func (r *AddOrderResponse) HasError() bool {
	return len(r.Error) > 0
}

func (c *RestClient) AddOrder(order AddOrderRequest) (*AddOrderResponse, error) {
	params := map[string]interface{}{}
	params["pair"] = order.Pair
	params["type"] = order.Side
	params["ordertype"] = order.Type
	params["price"] = fmt.Sprintf("%.8f", order.Price)
	params["volume"] = fmt.Sprintf("%.8f", order.Volume)
	if order.UserRef > 0 {
		params["userref"] = order.UserRef
	}
	if order.ValidateOnly {
		params["validate"] = "1"
	}

	httpResponse, err := c.Post("/0/private/AddOrder", params)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response AddOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed unmarshal response: %v: %s",
			err, string(body))
	}

	return &response, nil
}
