package main

import (
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
	env.sess, err = mgo.Dial("localhost")
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
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		var prices []*price
		coll.Find(nil).Limit(100).Iter().All(&prices)
		c.JSON(200, gin.H{
			"message": prices,
		})
	})

	r.Run()
}
