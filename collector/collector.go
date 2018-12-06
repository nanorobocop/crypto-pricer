package main

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nanorobocop/crypto-pricer/ticker"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

// Env contains environment variables
type Env struct {
	sess     *mgo.Session
	resultCh chan *ticker.Result
}

var env = &Env{}

func (env *Env) initialize() (err error) {

	// logger
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// db
	env.sess, err = mgo.Dial(os.Getenv("MONGO"))
	if err != nil {
		log.Fatalf("Cannot initiate session with DB: %+v", err)
	}
	env.sess.SetMode(mgo.Monotonic, true)

	// results
	env.resultCh = make(chan *ticker.Result)

	return nil
}

func (env *Env) handleExchange(exchange ticker.Ticker) {
	ex := exchange.Get()
	for _, pair := range ex.Pairs {
		log.WithFields(log.Fields{"exchange": ex.Exchange, "pair": pair}).Infof("Fetching price")
		go func(pair string) {
			exchange.FetchPair(pair, env.resultCh)
		}(pair)
	}

}

func (env *Env) handleExchanges() {
	oneMinDuration, _ := time.ParseDuration("1m")

	var exchanges []ticker.Ticker
	exchanges = append(exchanges, &ticker.Bitflyer{
		E: ticker.Exchange{
			Exchange:    "bitflyer",
			Endpoint:    "https://api.bitflyer.com/v1/ticker?product_code=",
			Frequency:   oneMinDuration,
			FrequencyCh: *time.NewTicker(oneMinDuration),
			Pairs:       []string{"BTC_JPY", "ETH_BTC"},
		},
	})
	exchanges = append(exchanges, &ticker.Zaif{
		E: ticker.Exchange{
			Exchange:    "zaif",
			Endpoint:    "https://api.zaif.jp/api/1/ticker/",
			Frequency:   oneMinDuration,
			FrequencyCh: *time.NewTicker(oneMinDuration),
			Pairs:       []string{"btc_jpy", "eth_jpy", "eth_btc"},
		},
	})

	// fetch each pair for each exchange
	for _, exchange := range exchanges {
		go func(exchange ticker.Ticker) {
			for {
				select {
				case <-exchange.Get().FrequencyCh.C:
					go env.handleExchange(exchange)
				}
			}
		}(exchange)
	}
}

func (env *Env) commitResults() {
	coll := env.sess.DB("cryptopricer").C("prices")

	for {
		select {
		case result := <-env.resultCh:
			log.WithFields(log.Fields{"exchange": result.Exchange.Exchange, "pair": result.PairInfo.Pair, "price": result.PairInfo.Price, "timestamp": result.PairInfo.Timestamp, "error": result.Error}).Info("Commiting result")
			if result.Error != nil {
				continue
			}
			dbEntry := struct {
				Datetime time.Time
				Pair     string
				Price    float64
				Exchange string
			}{
				Datetime: result.PairInfo.Timestamp,
				Pair:     strings.ToLower(result.PairInfo.Pair),
				Price:    result.PairInfo.Price,
				Exchange: result.Exchange.Exchange,
			}
			coll.Insert(&dbEntry)
		}
	}
}

func main() {
	env.initialize()
	defer env.sess.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go env.commitResults()
	env.handleExchanges()
	wg.Wait()
}
