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
	"io/ioutil"
	"strings"
)

type AssetPairService struct {
	cacheFilename string
	pairs         map[string]*AssetPairInfo
}

func NewAssetPairService() *AssetPairService {
	return &AssetPairService{}
}

func (s *AssetPairService) SetCacheFilename(filename string) {
	s.cacheFilename = filename
}

func (s *AssetPairService) Refresh() error {
	client, err := NewRestClient("", "")
	if err != nil {
		return err
	}
	response, err := client.Get("/0/public/AssetPairs")
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if s.cacheFilename != "" {
		ioutil.WriteFile(s.cacheFilename, body, 0644)
	}
	var assetPairResponse AssetPairResponse
	if err := json.Unmarshal(body, &assetPairResponse); err != nil {
		return err
	}
	s.LoadResponse(assetPairResponse)
	return nil
}

func (s *AssetPairService) LoadResponse(response AssetPairResponse) {
	keys := []string{}
	for pair := range response.Result {
		response.Result[pair].Pair = pair
		keys = append(keys, pair)
	}
	for _, key := range keys {
		info := response.Result[key]
		altname := response.Result[key].AltName
		wsname := response.Result[key].WsName
		response.Result[altname] = info
		response.Result[wsname] = info
	}
	s.pairs = response.Result
}

func (s *AssetPairService) GetRestPair(pair string) string {
	info, ok := s.pairs[strings.ToUpper(pair)]
	if ok {
		return info.Pair
	}
	return ""
}

func (s *AssetPairService) GetWsPair(pair string) string {
	info, ok := s.pairs[strings.ToUpper(pair)]
	if ok {
		return info.WsName
	}
	return ""
}

type AssetPairResponse struct {
	Error  []interface{}            `json:"error"`
	Result map[string]*AssetPairInfo `json:"result"`
}

type AssetPairInfo struct {
	Pair    string // Not in JSON response.
	AltName string `json:"altname"`
	WsName  string `json:"wsname"`
}
