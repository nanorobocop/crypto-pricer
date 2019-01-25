package ticker

import "time"

// https://www.bitstamp.net/api/v2/ticker/{currency_pair}/
// Supported values for currency_pair: btcusd, btceur, eurusd, xrpusd, xrpeur, xrpbtc, ltcusd, ltceur, ltcbtc, ethusd, etheur, ethbtc, bchusd, bcheur, bchbtc
// https://www.bitstamp.net/api/v2/ticker/btcusd/
// {"high": "4355.08", "last": "4227.90", "timestamp": "1543448821", "bid": "4225.20", "vwap": "4098.80", "volume": "22935.53196751", "low": "3754.01", "ask": "4228.29", "open": "3771.01"}

// Bitstamp contains struct for www.bitstamp.net
type Bitstamp struct {
	Exchange    string
	Endpoint    string
	Frequency   time.Duration
	Pairs       map[string]*PairInfo
	APIResponse struct {
		Price     string `json:"last"`
		Timestamp string `json:"timestamp"` // unix timestamp
	}
}
