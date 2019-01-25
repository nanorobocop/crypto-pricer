package ticker

// func TestGetPair(t *testing.T) {

// 	steps := []struct {
// 		json      string
// 		pair      string
// 		err       error
// 		last      float64
// 		timestamp time.Time
// 	}{
// 		{
// 			json:      `{"high": "4355.08", "last": "4227.90", "timestamp": "1543448821", "bid": "4225.20", "vwap": "4098.80", "volume": "22935.53196751", "low": "3754.01", "ask": "4228.29", "open": "3771.01"}`,
// 			pair:      "btcusd",
// 			err:       nil,
// 			last:      4227.90,
// 			timestamp: time.Unix(1543448821, 0),
// 		},
// 	}

// 	for i, step := range steps {
// 		t.Logf("[TEST] %d, %+v", i, step)

// 		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			fmt.Fprintln(w, step.json)
// 		}))
// 		defer ts.Close()

// 		pairs := make(map[string]*PairInfo)
// 		pairs["btcusd"] = &PairInfo{}
// 		ticker := Ticker{
// 			Exchange: "bitstamp.net",
// 			Endpoint: ts.URL,
// 			Pairs:    pairs,
// 		}

// 		actualErr := ticker.GetPair(step.pair)

// 		if (actualErr != nil && step.err != nil) || (ticker.Pairs[step.pair].Price == step.last && ticker.Pairs[step.pair].Timestamp == step.timestamp) {
// 			t.Logf("[TEST PASSED] %+v", ticker)
// 		} else {
// 			t.Errorf("[TEST FAILED] %+v, %+v", ticker, actualErr)
// 		}
// 	}
// }
