package ticker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// https://api.zaif.jp/api/1/ticker/btc_jpy
// {"last": 430645.0, "high": 457000.0, "low": 426015.0, "vwap": 440014.7346, "volume": 847.9854, "bid": 430645.0, "ask": 431075.0}

// Zaif contains struct for zaif.jp
type Zaif struct {
	E Exchange
}

// ZaifAPIResponse describes response
type ZaifAPIResponse struct {
	Price float64 `json:"last"`
}

// FetchPair returns rate for pair
func (t *Zaif) FetchPair(pair string, resultCh chan *Result) {
	u, err := url.Parse(t.E.Endpoint)
	// u.Path = u.Path + pair
	resp, err := http.Get(u.String() + pair)
	if err != nil {
		resultCh <- &Result{Error: err}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultCh <- &Result{Error: err}
		return
	}

	data := ZaifAPIResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		resultCh <- &Result{Error: err}
		return
	}

	resultCh <- &Result{Error: nil, Exchange: t.E, PairInfo: PairInfo{Pair: pair, Price: data.Price, Timestamp: time.Now()}}
}

// Get returns exchange
func (t *Zaif) Get() Exchange {
	return t.E
}
