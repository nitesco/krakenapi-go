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

type TimeResponse struct {
	Error  []interface{} `json:"error"`
	Result struct {
		UnixTime int64  `json:"unixtime"`
		Rfc1123  string `json:"rfc1123"`
	} `json:"result"`
}

func (c *RestClient) Time() (*TimeResponse, error) {
	httpReponse, err := c.Get("/0/public/Time")
	if err != nil {
		return nil, err
	}
	var response TimeResponse
	if err := decodeHttpResponse(httpReponse, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type CancelOrderResponse struct {
	Error  []interface{} `json:"error"`
	Result struct {
		Count int64 `json:"count"`
	} `json:"result"`
}

func (r *CancelOrderResponse) HasError() bool {
	return len(r.Error) > 0
}

func (c *RestClient) CancelOrder(txId string) (*CancelOrderResponse, error) {
	params := map[string]interface{}{
		"txid": txId,
	}
	httpResponse, err := c.Post("/0/private/CancelOrder", params)
	if err != nil {
		return nil, err
	}
	var response CancelOrderResponse
	if err := decodeHttpResponse(httpResponse, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
