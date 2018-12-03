package main

import (
	"sync"
	"time"

	"github.com/nanorobocop/crypto-pricer/ticker"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

// Env contains environment variables
type Env struct {
	sess *mgo.Session
}

var env = &Env{}

func (env *Env) initialize() (err error) {

	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	env.sess, err = mgo.Dial("localhost")
	if err != nil {
		log.Error("Cannot initiate session with DB")
	}

	env.sess.SetMode(mgo.Monotonic, true)

	return nil
}

func main() {
	env.initialize()
	defer env.sess.Close()

	frequency, err := time.ParseDuration("5m")
	if err != nil {
		log.Fatal(err)
	}

	var exchanges []*ticker.Ticker
	exchanges = append(exchanges, ticker.NewTickerAPI("bitstamp.net", "https://www.bitstamp.net/api/v2/ticker", frequency, []string{"btcusd", "btceur"}))

	coll := env.sess.DB("cryptopricer").C("prices")

	var wg sync.WaitGroup
	for _, exchange := range exchanges {
		wg.Add(1)
		log.WithFields(log.Fields{"exchange": exchange.Exchange}).Infof("Requesting info from exchange: %+v", exchange)
		go func(exchange *ticker.Ticker) {
			defer wg.Done()
			exchange.GetAllPairs()
		}(exchange)
	}
	wg.Wait()

	for _, exchange := range exchanges {
		for pairName, pairInfo := range exchange.Pairs {
			log.WithFields(log.Fields{"exchange": exchange.Exchange, "pair": pairName, "price": pairInfo.Price, "timestamp": pairInfo.Timestamp}).Info("Saving data")
			dbEntry := struct {
				Datetime time.Time
				Pair     string
				Price    float64
				Exchange string
			}{
				Datetime: pairInfo.Timestamp,
				Pair:     pairName,
				Price:    pairInfo.Price,
				Exchange: exchange.Exchange,
			}
			coll.Insert(&dbEntry)
		}
	}
}
