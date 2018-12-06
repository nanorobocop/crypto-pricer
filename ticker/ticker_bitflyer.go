package ticker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// https://api.bitflyer.com/v1/ticker?product_code=BTC_JPY
// {"product_code":"BTC_JPY","timestamp":"2018-12-04T04:54:58.453","tick_id":6405348,"best_bid":431801.0,"best_ask":432131.0,"best_bid_size":1.50000026,"best_ask_size":1.57298296,"total_bid_depth":2337.77697571,"total_ask_depth":2198.29047712,"ltp":432154.0,"volume":473281.98687278,"volume_by_product":10745.03169861}

// Bitflyer contains struct for bitflyer.com
type Bitflyer struct {
	E Exchange
}

// BitflyerAPIResponse describes response
type BitflyerAPIResponse struct {
	Price     float64 `json:"best_bid"`
	Timestamp string  `json:"timestamp"` // ISO format
}

// FetchPair returns rate for pair
func (t *Bitflyer) FetchPair(pair string, resultCh chan *Result) {

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

	data := BitflyerAPIResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		resultCh <- &Result{Error: err}
		return
	}

	data.Timestamp = data.Timestamp[0:19]
	timestamp, err := time.ParseInLocation("2006-01-02T15:04:05", data.Timestamp, &time.Location{})
	if err != nil {
		resultCh <- &Result{Error: err}
		return
	}

	resultCh <- &Result{Error: nil, Exchange: t.E, PairInfo: PairInfo{Pair: pair, Price: data.Price, Timestamp: timestamp}}
}

// Get returns exchange
func (t *Bitflyer) Get() Exchange {
	return t.E
}
