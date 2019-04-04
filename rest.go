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
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var API_ROOT string = "https://api.kraken.com"

func init() {
	if os.Getenv("KRAKEN_API_ROOT") != "" {
		API_ROOT = os.Getenv("KRAKEN_API_ROOT")
	}
}

type RestClient struct {
	apiKey    string
	apiSecret []byte
	lastNonce int64
	lock      sync.Mutex
}

func NewRestClient(apiKey string, apiSecret string) (*RestClient, error) {
	var decodedApiSecret []byte = nil
	var err error
	if apiKey != "" && apiSecret != "" {
		decodedApiSecret, err = base64.StdEncoding.DecodeString(apiSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to base64 decode api secret: %v", err)
		}
	}

	return &RestClient{
		apiKey:    apiKey,
		apiSecret: decodedApiSecret,
	}, nil
}

func (c *RestClient) Get(path string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", API_ROOT, strings.TrimPrefix(path, "/"))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(request)
}

func (c *RestClient) Post(path string, params map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", API_ROOT, strings.TrimPrefix(path, "/"))
	if params == nil {
		params = map[string]interface{}{}
	}
	nonce := c.getNonce()
	params["nonce"] = nonce
	queryString := c.buildQueryString(params)
	request, err := http.NewRequest("POST", url, strings.NewReader(queryString))
	if err != nil {
		return nil, err
	}
	c.authenticateRequest(request, path, nonce, queryString)
	return http.DefaultClient.Do(request)
}

func (c *RestClient) getNonce() int64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	nonce := time.Now().UnixNano() / int64(time.Millisecond)
	if nonce == c.lastNonce {
		nonce += 1
	}
	c.lastNonce = nonce
	return nonce
}

func (c *RestClient) authenticateRequest(request *http.Request, endpoint string, nonce int64, postData string) {
	s256 := sha256.New()
	s256.Write([]byte(fmt.Sprintf("%d%s", nonce, postData)))

	mac := hmac.New(sha512.New, c.apiSecret)
	mac.Write([]byte(endpoint))
	mac.Write(s256.Sum(nil))

	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	request.Header.Add("API-Key", c.apiKey)
	request.Header.Add("API-Sign", signature)
}

func (c *RestClient) buildQueryString(params map[string]interface{}) string {
	queryString := ""

	keys := func() []string {
		keys := []string{}
		for key, _ := range params {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys
	}()

	for _, key := range keys {
		if queryString != "" {
			queryString = fmt.Sprintf("%s&", queryString)
		}
		queryString = fmt.Sprintf("%s%s=%v", queryString, key, params[key])
	}

	return queryString
}
