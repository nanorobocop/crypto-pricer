package ticker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBitflyerGetPair(t *testing.T) {

	steps := []struct {
		json      string
		pair      string
		err       error
		price     float64
		timestamp time.Time
	}{
		{
			json:      `{"product_code":"BTC_JPY","timestamp":"2018-12-04T04:54:58.453","tick_id":6405348,"best_bid":431801.0,"best_ask":432131.0,"best_bid_size":1.50000026,"best_ask_size":1.57298296,"total_bid_depth":2337.77697571,"total_ask_depth":2198.29047712,"ltp":432154.0,"volume":473281.98687278,"volume_by_product":10745.03169861}`,
			pair:      "BTC_JPY",
			err:       nil,
			price:     431801.0,
			timestamp: time.Unix(1543899298, 0),
		},
	}

	for i, step := range steps {
		t.Logf("[TEST] %d, %+v", i, step)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, step.json)
		}))
		defer ts.Close()

		ticker := Bitflyer{
			E: Exchange{
				Exchange: "bitstamp.net",
				Endpoint: ts.URL,
				Pairs:    []string{"BTC_JPY"},
			},
		}
		resultCh := make(chan *Result)
		ticker.FetchPair(step.pair, resultCh)

		// if (actualErr != nil && step.err != nil) || (ticker.E.Pairs[step.pair].Price == step.price && reflect.DeepEqual(ticker.E.Pairs[step.pair].Timestamp, step.timestamp)) {
		// 	t.Logf("[TEST PASSED] %+v", ticker)
		// } else {
		// 	t.Errorf("[TEST FAILED] %+v, %+v", ticker, actualErr)
		// }
	}
}
