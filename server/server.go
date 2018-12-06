package main

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Env contains environment variables
type Env struct {
	sess *mgo.Session
}

var env = &Env{}

func (env *Env) initialize() (err error) {
	env.sess, err = mgo.Dial(os.Getenv("MONGO"))
	if err != nil {
		log.Error("Cannot initiate session with DB")
	}

	env.sess.SetMode(mgo.Monotonic, true)

	return nil
}

type price struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Exchange string        `json:"exchange" bson:"exchange"`
	Pair     string        `json:"pair" bson:"pair"`
	Price    float64       `json:"price" bson:"price"`
	Datetime time.Time     `json:"datetime" bson:"datetime"`
}

func main() {
	env.initialize()
	defer env.sess.Close()

	coll := env.sess.DB("cryptopricer").C("prices")

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		var prices []*price
		coll.Find(nil).Limit(100).Iter().All(&prices)
		c.JSON(200, gin.H{
			"message": prices,
		})
	})

	r.GET("/alerts", func(c *gin.Context) {
		type minmaxStruct struct {
			ID       string    `json:"pair" bson:"_id"`
			Min      float64   `json:"min" bson:"min"`
			Max      float64   `json:"max" bson:"max"`
			MaxMin   float64   `json:"price_between_max_min" bson:"diff"`
			Datetime time.Time `json:"datetime" bson:"datetime"`
			Alert    bool      `json:"alert" bson:"alert"`
		}

		var datetime time.Time
		period, err := time.ParseDuration(c.DefaultQuery("period", "1h"))
		if err != nil {
			datetime = time.Now().Add(-time.Hour)
		} else {
			datetime = time.Now().Add(-period)
		}

		thresholdArg := c.DefaultQuery("threshold", "10")
		threshold, err := strconv.ParseFloat(thresholdArg, 64)
		if err != nil {
			thresholdArg = "10"
		}

		// db.prices.aggregate([{$match: {"datetime" : { $gte : new Date(ISODate().getTime() - 1000*60*280) }}},{ $group: {_id: "$pair", min: { $min: "$price"}, max: {$max: "$price" } }}, { $project: {div: {$divide: [5, 4]}, min: "$min", max: "$max", div: {$divide: ["$max", "$min"]}  }}]);
		// { "_id" : "btc_jpy", "div" : 1.003413891004778, "min" : 433230, "max" : 434709 }

		var minmax []*minmaxStruct
		pipe := coll.Pipe([]bson.M{
			{
				"$match": bson.M{
					"datetime": bson.M{
						"$gte": datetime,
					},
				},
			},
			{
				"$group": bson.M{
					"_id": "$pair",
					"min": bson.M{"$min": "$price"},
					"max": bson.M{"$max": "$price"},
				},
			}, {
				"$project": bson.M{
					"diff": bson.M{
						"$subtract": []string{"$max", "$min"},
					},
					"min":      "$min",
					"max":      "$max",
					"datetime": datetime,
					"alert": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$gte": []interface{}{
									bson.M{
										"$subtract": []string{"$max", "$min"},
									},
									threshold,
								},
							},
							"then": true,
							"else": false,
						},
					},
				},
			},
		})
		pipe.Iter().All(&minmax)
		c.JSON(200, gin.H{
			"message": minmax,
		})
	})
	r.Run()
}
