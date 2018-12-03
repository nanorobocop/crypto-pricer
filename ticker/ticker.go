package ticker

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"
)

// https://www.bitstamp.net/api/v2/ticker/{currency_pair}/
// Supported values for currency_pair: btcusd, btceur, eurusd, xrpusd, xrpeur, xrpbtc, ltcusd, ltceur, ltcbtc, ethusd, etheur, ethbtc, bchusd, bcheur, bchbtc
// https://www.bitstamp.net/api/v2/ticker/btcusd/
// {"high": "4355.08", "last": "4227.90", "timestamp": "1543448821", "bid": "4225.20", "vwap": "4098.80", "volume": "22935.53196751", "low": "3754.01", "ask": "4228.29", "open": "3771.01"}

// Ticker contains info about ticker API
type Ticker struct {
	Exchange  string
	Endpoint  string
	Frequency time.Duration
	Pairs     map[string]*PairInfo
}

// PairInfo contains info about pair
type PairInfo struct {
	Price     float64
	Timestamp time.Time
}

// APIResponse is structure for response
type APIResponse struct {
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	Bid       string `json:"bid"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}

// NewTickerAPI creates Ticker struct
func NewTickerAPI(exchange, endpoint string, frequency time.Duration, pairs []string) *Ticker {
	pp := make(map[string]*PairInfo)
	for _, pair := range pairs {
		pp[pair] = &PairInfo{}
	}
	return &Ticker{Exchange: exchange, Endpoint: endpoint, Frequency: frequency, Pairs: pp}
}

// GetPair returns rate for pair
func (t *Ticker) GetPair(pair string) (err error) {
	if _, ok := t.Pairs[pair]; !ok {
		return errors.New("pair is not a key of Pairs")
	}

	u, err := url.Parse(t.Endpoint)
	u.Path = path.Join(u.Path, pair)
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respData APIResponse
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	timestamp, err := strconv.ParseInt(respData.Timestamp, 10, 64)
	if err != nil {
		return err
	}

	last, err := strconv.ParseFloat(respData.Last, 64)

	t.Pairs[pair] = &PairInfo{
		Price:     last,
		Timestamp: time.Unix(timestamp, 0),
	}
	return nil
}

// GetAllPairs works like GetPair but for all pairs
func (t *Ticker) GetAllPairs() {
	var wg sync.WaitGroup
	for key := range t.Pairs {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			t.GetPair(key)
		}(key)
	}
	wg.Wait()
}
