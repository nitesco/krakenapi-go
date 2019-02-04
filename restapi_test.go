package krakenapi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddOrderResponseDecode(t *testing.T) {
	responseBody := `{"error":[],"result":{"descr":{"order":"buy 0.01000000 XBTUSD @ limit 3000.0"},"txid":["OF4HUU-ZBMIF-BG2R4Z"]}}`
	var response AddOrderResponse
	err := json.Unmarshal([]byte(responseBody), &response)
	assert.Nil(t, err)
	assert.Len(t, response.Error, 0)
	assert.Equal(t, response.Result.Descr.Order, "buy 0.01000000 XBTUSD @ limit 3000.0")
	assert.Len(t, response.Result.Txid, 1)
	assert.Equal(t, response.Result.Txid[0], "OF4HUU-ZBMIF-BG2R4Z")
}
