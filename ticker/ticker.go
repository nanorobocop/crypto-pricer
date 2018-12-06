package ticker

import (
	"time"
)

// Ticker provides interface for Ticker API
type Ticker interface {
	FetchPair(string, chan *Result)
	Get() Exchange
}

// Exchange contains ticker information
type Exchange struct {
	Exchange    string
	Endpoint    string
	Frequency   time.Duration
	FrequencyCh time.Ticker
	Pairs       []string
}

// PairInfo contains info about pair
type PairInfo struct {
	Pair      string
	Price     float64
	Timestamp time.Time
}

// Result contains results
type Result struct {
	Error    error
	Exchange Exchange
	PairInfo PairInfo
}
